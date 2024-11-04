package service

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/config"
	"kbox-api/internal/model"
	"kbox-api/internal/modules/auth/dto"
	"kbox-api/internal/modules/user/handler"
	"kbox-api/pkg/token"
	"kbox-api/shared/httperror"
)

type AuthService struct {
	userService handler.UserService
	config      *config.Config
}

func NewAuthService(userService handler.UserService, config *config.Config) *AuthService {
	return &AuthService{
		userService: userService,
		config:      config,
	}
}

func (s *AuthService) Login(dto dto.LoginDTO) (string, error) {
	user, err := s.userService.GetOneByEmail(dto.Email)
	if err != nil || !s.userService.CheckPassword(user, dto.Password) {
		return "", httperror.New(
			http.StatusUnauthorized,
			"Неверный логин или пароль",
		)
	}
	jwtKey := []byte(s.config.App.JwtSecretKey)
	expiresIn, err := time.ParseDuration(s.config.App.JwtExpiresIn)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	tokenString, err := token.GenerateToken(jwtKey, user.Id, user.Role, expiresIn)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError, err.Error(),
		)
	}
	return tokenString, nil
}

func (s *AuthService) Register(user dto.RegisterDTO) (string, error) {
	existingUser, err := s.userService.GetOneByEmail(user.Email)
	if err == nil && existingUser.Id != "" {
		return "", httperror.New(
			http.StatusConflict,
			"Пользователь с таким email уже зарегистрирован",
		)
	}
	existingUser, err = s.userService.GetOneByNickName(user.NickName)
	if err == nil && existingUser.Id != "" {
		return "", httperror.New(
			http.StatusConflict,
			"Пользователь с таким ником уже зарегистрирован",
		)
	}
	hashedPassword, err := s.userService.HashPassword(user.Password)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	newUser := model.User{
		NickName:   user.NickName,
		Name:       user.Name,
		Email:      user.Email,
		Password:   hashedPassword,
		Photo:      user.Photo,
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

	jwtKey := []byte(s.config.App.JwtSecretKey)

	tokenString, err := token.GenerateToken(jwtKey, createdUser.Id, createdUser.Role, 24*time.Hour)
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
