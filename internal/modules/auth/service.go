package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/model"
	externalservice "github.com/1Storm3/flibox-api/internal/modules/external"
	dtoUser "github.com/1Storm3/flibox-api/internal/modules/user"
	"github.com/1Storm3/flibox-api/internal/shared/helper"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"github.com/1Storm3/flibox-api/pkg/token"
	"go.uber.org/zap"
)

var _ ServiceInterface = (*Service)(nil)

type ServiceInterface interface {
	Login(ctx context.Context, dto LoginDTO) (string, error)
	Register(ctx context.Context, user RegisterDTO) (bool, error)
	Me(ctx context.Context, userId string) (MeResponseDTO, error)
	Verify(ctx context.Context, tokenDto string) error
}

type Service struct {
	userService  dtoUser.ServiceInterface
	emailService externalservice.EmailServiceInterface
	cfg          *config.Config
}

func NewAuthService(userService dtoUser.ServiceInterface,
	emailService externalservice.EmailServiceInterface,
	cfg *config.Config,
) *Service {
	return &Service{
		userService:  userService,
		emailService: emailService,
		cfg:          cfg,
	}
}

func (s *Service) Login(ctx context.Context, req LoginDTO) (string, error) {
	user, err := s.userService.GetOneByEmail(ctx, req.Email)
	if err != nil || !s.userService.CheckPassword(ctx, &user, req.Password) {
		return "", httperror.New(
			http.StatusUnauthorized,
			"Неверный логин или пароль",
		)
	}
	jwtKey := []byte(s.cfg.App.JwtSecretKey)
	expiresIn, err := time.ParseDuration(s.cfg.App.JwtExpiresIn)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	tokenString, err := token.GenerateToken(jwtKey, user.ID, user.Role, expiresIn)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError, err.Error(),
		)
	}
	return tokenString, nil
}

func (s *Service) Register(ctx context.Context, req RegisterDTO) (bool, error) {
	existingUser, err := s.userService.GetOneByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != "" {
		return false, httperror.New(
			http.StatusConflict,
			"Пользователь с таким email уже зарегистрирован",
		)
	}
	existingUser, err = s.userService.GetOneByNickName(ctx, req.NickName)
	if err == nil && existingUser.ID != "" {
		return false, httperror.New(
			http.StatusConflict,
			"Пользователь с таким ником уже зарегистрирован",
		)
	}
	hashedPassword, err := s.userService.HashPassword(ctx, req.Password)
	if err != nil {
		return false, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	jwtKey := []byte(s.cfg.App.JwtSecretKey)

	verificationToken, err := token.GenerateEmailToken(req.Email, jwtKey, time.Hour*2)
	if err != nil {
		return false, httperror.New(
			http.StatusInternalServerError,
			"Не удалось создать токен для подтверждения email",
		)
	}

	newUser := model.User{
		NickName:      req.NickName,
		Name:          req.Name,
		Email:         req.Email,
		Password:      hashedPassword,
		Photo:         req.Photo,
		Role:          "user",
		IsVerified:    false,
		VerifiedToken: verificationToken,
		LastActivity:  time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	createdUser, err := s.userService.Create(ctx, newUser)
	if err != nil {
		return false, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	go func() {
		emailBody, err := helper.TakeHTMLTemplate(s.cfg.App.AppUrl, *verificationToken)
		if err != nil {
			logger.Error(err.Error())
		}
		if err := s.emailService.SendEmail(
			createdUser.Email,
			emailBody,
			"Подтверждение регистрации",
		); err != nil {
			logger.Error("Ошибка при отправке email", zap.Error(err))
		}
	}()

	return true, nil
}

func (s *Service) Verify(ctx context.Context, tokenDto string) error {
	jwtKey := []byte(s.cfg.App.JwtSecretKey)
	email, err := token.ValidateEmailToken(tokenDto, jwtKey)
	if err != nil {
		return httperror.New(
			http.StatusBadRequest,
			"Неверный токен",
		)
	}
	userNotVerified, err := s.userService.GetOneByEmail(ctx, email)
	if err != nil {
		return httperror.New(
			http.StatusNotFound,
			"Пользователь не найден",
		)
	}
	user := dtoUser.UpdateForVerifyDTO{
		ID:            userNotVerified.ID,
		IsVerified:    true,
		VerifiedToken: nil,
	}

	if _, err := s.userService.UpdateForVerify(ctx, user); err != nil {
		return httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return nil
}

func (s *Service) Me(ctx context.Context, userId string) (MeResponseDTO, error) {

	user, err := s.userService.GetOneById(ctx, userId)

	if err != nil {
		return MeResponseDTO{}, httperror.New(
			http.StatusNotFound,
			"Пользователь не найден",
		)
	}

	lastActivity := time.Now().Format(time.RFC3339)
	_, err = s.userService.Update(ctx, dtoUser.UpdateUserDTO{
		ID:           userId,
		LastActivity: &lastActivity,
	})

	if err != nil {
		return MeResponseDTO{}, httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return MapModelUserToResponseDTO(user), nil
}
