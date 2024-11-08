package service

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/config"
	"kbox-api/internal/model"
	"kbox-api/internal/modules/auth/dto"
	"kbox-api/internal/modules/user/service"
	"kbox-api/pkg/token"
	"kbox-api/shared/httperror"
)

type AuthServiceInterface interface {
	Login(dto dto.LoginDTO) (string, error)
	Register(user dto.RegisterDTO) (string, error)
	Me(c *fiber.Ctx) (model.User, error)
}

type AuthService struct {
	userService service.UserServiceInterface
	cfg         *config.Config
}

func NewAuthService(userService service.UserServiceInterface, cfg *config.Config) *AuthService {
	return &AuthService{
		userService: userService,
		cfg:         cfg,
	}
}

func (s *AuthService) Login(req dto.LoginDTO) (string, error) {
	user, err := s.userService.GetOneByEmail(req.Email)
	if err != nil || !s.userService.CheckPassword(user, req.Password) {
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

func (s *AuthService) Register(req dto.RegisterDTO) (string, error) {
	existingUser, err := s.userService.GetOneByEmail(req.Email)
	if err == nil && existingUser.ID != "" {
		return "", httperror.New(
			http.StatusConflict,
			"Пользователь с таким email уже зарегистрирован",
		)
	}
	existingUser, err = s.userService.GetOneByNickName(req.NickName)
	if err == nil && existingUser.ID != "" {
		return "", httperror.New(
			http.StatusConflict,
			"Пользователь с таким ником уже зарегистрирован",
		)
	}
	hashedPassword, err := s.userService.HashPassword(req.Password)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	newUser := model.User{
		NickName:   req.NickName,
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Photo:      req.Photo,
		Role:       "user",
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdUser, err := s.userService.Create(newUser)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	jwtKey := []byte(s.cfg.App.JwtSecretKey)

	tokenString, err := token.GenerateToken(jwtKey, createdUser.ID, createdUser.Role, 24*time.Hour)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return tokenString, nil
}

func (s *AuthService) Me(c *fiber.Ctx) (model.User, error) {
	claims, ok := c.Locals("userClaims").(*token.Claims)
	if !ok {
		return model.User{}, httperror.New(
			http.StatusUnauthorized,
			"Не удалось получить информацию о пользователе",
		)
	}

	user, err := s.userService.GetOneById(claims.UserID)
	if err != nil {
		return model.User{}, httperror.New(
			http.StatusNotFound,
			"Пользователь не найден",
		)
	}

	return user, nil
}
