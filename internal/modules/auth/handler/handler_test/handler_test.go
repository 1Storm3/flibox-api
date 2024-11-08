package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"kbox-api/internal/delivery/middleware"
	"kbox-api/internal/model"
	"kbox-api/internal/modules/auth/dto"
	"kbox-api/internal/modules/auth/handler"
	"kbox-api/internal/modules/auth/service"
)

var jwtKey = "778c4001-18d0-4e13-8699-6b398f196b1a"

var jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzAxMzRjNWMtM2UwNi00NjBjLWFmZjctOTg3NDBmZDNhOTA0Iiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzEwMDY2MTB9.VDw7oB0mesaEaFTnN-qxRzXXUWPB4W8xQ1KgkQEIM-o"

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(loginDTO dto.LoginDTO) (string, error) {
	args := m.Called(loginDTO)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Register(registerDTO dto.RegisterDTO) (string, error) {
	args := m.Called(registerDTO)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Me(c *fiber.Ctx) (model.User, error) {
	args := m.Called(c)
	return args.Get(0).(model.User), args.Error(1)
}

func setupTestApp(authService service.AuthServiceInterface) *fiber.App {
	app := fiber.New()
	authHandler := handler.NewAuthHandler(authService)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("jwtKey", jwtKey)
		return c.Next()
	})

	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/register", authHandler.Register)
	app.Get("/auth/me", middleware.AuthMiddleware, authHandler.Me)

	return app
}

func TestLogin(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Login", dto.LoginDTO{Email: "user@example.com", Password: "password123"}).Return("mock_token", nil)

	app := setupTestApp(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"user@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "mock_token", body["token"])

	mockAuthService.AssertExpectations(t)
}

func TestRegister(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Register", dto.RegisterDTO{Email: "newuser@example.com", Password: "password123"}).Return("mock_token", nil)

	app := setupTestApp(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(
		`{"email":"newuser@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "mock_token", body["token"])

	mockAuthService.AssertExpectations(t)
}

func TestMe(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Me", mock.Anything).Return(
		model.User{Id: "70134c5c-3e06-460c-aff7-98740fd3a904", Email: "user@example.com"}, nil)

	app := setupTestApp(mockAuthService)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", jwtToken)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var userResponse map[string]interface{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&userResponse))
	assert.Equal(t, "70134c5c-3e06-460c-aff7-98740fd3a904", userResponse["id"])
	assert.Equal(t, "user@example.com", userResponse["email"])

	mockAuthService.AssertExpectations(t)
}
