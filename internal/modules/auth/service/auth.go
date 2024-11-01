package service

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/modules/user/handler"
	"kinopoisk-api/internal/modules/user/service"
	"kinopoisk-api/pkg/token"
	"kinopoisk-api/shared/httperror"
	"net/http"
	"time"
)

type AuthService struct {
	userService handler.UserService
}

type RequestUser struct {
	Name     string `json:"name" validate:"required,min=4"`
	NickName string `json:"nickName" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Photo    string `json:"photo" validate:"omitempty,url"`
}

func NewAuthService(userService handler.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Login(login string, password string) (string, error) {
	user, err := s.userService.GetOneByEmail(login)
	if err != nil || !s.userService.CheckPassword(user, password) {
		return "", httperror.New(
			http.StatusUnauthorized,
			"Неверный логин или пароль",
		)
	}
	tokenString, err := token.GenerateToken(user.Id, user.Role, time.Hour*24)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError, err.Error(),
		)
	}
	return tokenString, nil
}

func (s *AuthService) Register(user RequestUser) (string, error) {
	existingUser, err := s.userService.GetOneByEmail(user.Email)
	if err == nil && existingUser.Id != "" {
		return "", httperror.New(
			http.StatusConflict,
			"Пользователь с таким email уже зарегистрирован",
		)
	}
	hashedPassword, err := s.userService.HashPassword(user.Password)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	newUser := service.User{
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

	createdUser, err := s.userService.CreateUser(newUser)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	tokenString, err := token.GenerateToken(createdUser.Id, createdUser.Role, 24*time.Hour)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return tokenString, nil
}

func (s *AuthService) Me(c *fiber.Ctx) (service.User, error) {
	claims, ok := c.Locals("userClaims").(*token.Claims)
	if !ok {
		return service.User{}, httperror.New(
			http.StatusUnauthorized,
			"Не удалось получить информацию о пользователе",
		)
	}

	user, err := s.userService.GetOneById(claims.UserID)
	if err != nil {
		return service.User{}, httperror.New(
			http.StatusNotFound,
			"Пользователь не найден",
		)
	}

	return user, nil
}
