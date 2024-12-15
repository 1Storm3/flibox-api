package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/1Storm3/flibox-api/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testUserLoginDTO = LoginDTO{
		Email:    "user@example.com",
		Password: "password123",
	}

	testUserRegisterDTO = RegisterDTO{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "John Doe",
		NickName: "johndoe",
	}

	testUserMe = MeResponseDTO{
		Id:    "cd1d4e27-f774-4d6b-ae8d-9ef58ec79711",
		Email: "user@example.com",
	}
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Verify(ctx context.Context, tokenDto string) error {
	args := m.Called(tokenDto)
	return args.Error(0)
}

func (m *MockAuthService) Login(ctx context.Context, dto LoginDTO) (string, error) {
	args := m.Called(dto)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, user RegisterDTO) (bool, error) {
	args := m.Called(user)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) Me(ctx context.Context, userId string) (MeResponseDTO, error) {
	args := m.Called(userId)
	return args.Get(0).(MeResponseDTO), args.Error(1)
}

func TestVerify(t *testing.T) {
	t.Run("Verify success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		mockAuthService.On("Verify", "mock_token").Return(nil)

		app := fiber.New()
		authHandler := NewAuthHandler(mockAuthService)

		app.Post("/auth/verify/:token", authHandler.Verify)

		req := httptest.NewRequest(http.MethodPost, "/auth/verify/mock_token", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		expectedResponse := `{"message":"Пользователь верифицирован"}`
		assert.JSONEq(t, expectedResponse, string(body))

		mockAuthService.AssertExpectations(t)
	})

	t.Run("Verify failed", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		mockAuthService.On("Verify", "invalid_token").Return(errors.New("verification failed"))

		app := fiber.New()
		authHandler := NewAuthHandler(mockAuthService)

		app.Post("/auth/verify/:token", authHandler.Verify)

		req := httptest.NewRequest(http.MethodPost, "/auth/verify/invalid_token", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "verification failed")

		mockAuthService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Bad Request", func(t *testing.T) {
		authHandler := NewAuthHandler(nil)
		app := fiber.New()

		app.Post("/auth/login", authHandler.Login)

		reqInvalid := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email": "user@example.com", "password":34}`))
		reqInvalid.Header.Set("Content-Type", "application/json")

		respInvalidBadReq, err := app.Test(reqInvalid)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, respInvalidBadReq.StatusCode)

		bodyInvalidBadReq, _ := io.ReadAll(respInvalidBadReq.Body)

		var errorResponse map[string]interface{}

		assert.NoError(t, json.Unmarshal(bodyInvalidBadReq, &errorResponse))

		assert.Equal(t, "Некорректные данные запроса", errorResponse["message"])
		assert.Equal(t, float64(400), errorResponse["statusCode"])
	})

	t.Run("Login Success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		mockAuthService.On("Login", testUserLoginDTO).Return("mock_token", nil)

		authHandler := NewAuthHandler(mockAuthService)
		app := fiber.New()

		app.Post("/auth/login", authHandler.Login)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"user@example.com","password":"password123"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)

		var response map[string]interface{}
		_ = json.Unmarshal(body, &response)

		assert.NoError(t, err)
		assert.Equal(t, "mock_token", response["token"])

		mockAuthService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		unexpectedError := errors.New("Неожиданная ошибка")
		mockAuthService.On("Login", mock.Anything, mock.Anything).Return("", unexpectedError)

		authHandler := NewAuthHandler(mockAuthService)
		app := fiber.New()

		app.Post("/auth/login", authHandler.Login)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email": "test@example.com", "password": "1234"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		err := json.Unmarshal(body, &response)
		if err != nil {
			return
		}

		assert.Equal(t, "Неожиданная ошибка", response["message"])
		assert.Equal(t, float64(http.StatusInternalServerError), response["statusCode"])
		assert.NoError(t, err)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		mockAuthService.On("Login", LoginDTO{Email: "user@example.com", Password: "wrongpassword"}).
			Return("", httperror.New(http.StatusUnauthorized, "Неверный логин или пароль"))

		authHandler := NewAuthHandler(mockAuthService)
		app := fiber.New()

		app.Post("/auth/login", authHandler.Login)

		reqError := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"user@example.com","password":"wrongpassword"}`))
		reqError.Header.Set("Content-Type", "application/json")

		respErrorWrong, err := app.Test(reqError)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, respErrorWrong.StatusCode)

		var errorResponseWrong map[string]interface{}
		assert.NoError(t, json.NewDecoder(respErrorWrong.Body).Decode(&errorResponseWrong))

		assert.Equal(t, "Неверный логин или пароль", errorResponseWrong["message"])
		assert.Equal(t, float64(401), errorResponseWrong["statusCode"])

		mockAuthService.AssertExpectations(t)
	})
}
func TestRegister(t *testing.T) {
	t.Run("Bad Request", func(t *testing.T) {
		app := fiber.New()
		authHandler := NewAuthHandler(nil)

		app.Post("/auth/register", authHandler.Register)

		reqInvalid := httptest.
			NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(
				`{"email": "user@example.com", "password":34, "name": "John Doe", "nickName": "johndoe"}`))
		reqInvalid.Header.Set("Content-Type", "application/json")

		respInvalidBadReq, err := app.Test(reqInvalid)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, respInvalidBadReq.StatusCode)

		bodyInvalidBadReq, _ := io.ReadAll(respInvalidBadReq.Body)

		var errorResponse map[string]interface{}

		assert.NoError(t, json.Unmarshal(bodyInvalidBadReq, &errorResponse))

		assert.Equal(t, "Некорректные данные запроса", errorResponse["message"])
		assert.Equal(t, float64(400), errorResponse["statusCode"])
	})

	t.Run("Success register", func(t *testing.T) {
		mockAuthService := new(MockAuthService)

		mockAuthService.On("Register", testUserRegisterDTO).Return(true, nil)

		app := fiber.New()
		authHandler := NewAuthHandler(mockAuthService)

		app.Post("/auth/register", authHandler.Register)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(
			`{"email":"newuser@example.com","password":"password123","name": "John Doe", "nickName": "johndoe"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]bool
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, true, body["data"])

		mockAuthService.AssertExpectations(t)
	})

	t.Run("Failed register", func(t *testing.T) {
		mockAuthService := new(MockAuthService)

		mockAuthService.On("Register", RegisterDTO{
			Email:    "newuser@example.com",
			Password: "password123",
			Name:     "John Doe",
			NickName: "johndoe"},
		).Return(false,
			httperror.New(http.StatusInternalServerError, "failed register"))

		app := fiber.New()
		authHandler := NewAuthHandler(mockAuthService)

		app.Post("/auth/register", authHandler.Register)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(
			`{"email":"newuser@example.com","password":"password123","name": "John Doe", "nickName": "johndoe"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(body, &response))

		assert.Equal(t, "failed register", response["message"])
		assert.Equal(t, float64(http.StatusInternalServerError), response["statusCode"])

		mockAuthService.AssertExpectations(t)
	})
}

func TestMe(t *testing.T) {
	t.Run("Me success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)

		mockAuthService.On("Me", mock.Anything).Return(testUserMe, nil)

		app := fiber.New()
		authHandler := NewAuthHandler(mockAuthService)
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("userClaims", &token.Claims{UserID: "cd1d4e27-f774-4d6b-ae8d-9ef58ec79711"})
			return c.Next()
		})
		app.Get("/auth/me", authHandler.Me)

		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var userResponse MeResponseDTO
		assert.NoError(t, json.Unmarshal(body, &userResponse))

		assert.Equal(t, testUserMe.Id, userResponse.Id)
		assert.Equal(t, testUserMe.Email, userResponse.Email)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Me failed", func(t *testing.T) {
		mockAuthService := new(MockAuthService)

		mockAuthService.On("Me", mock.Anything).
			Return(MeResponseDTO{}, httperror.New(http.StatusUnauthorized, "Не удалось получить информацию о пользователе"))

		app := fiber.New()
		authHandler := NewAuthHandler(mockAuthService)
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("userClaims", &token.Claims{UserID: "cd1d4e27-f774-4d6b-ae8d-9ef58ec79711"})
			return c.Next()
		})
		app.Get("/auth/me", authHandler.Me)

		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var errorResponse map[string]interface{}
		assert.NoError(t, json.Unmarshal(body, &errorResponse))
		assert.Equal(t, "Не удалось получить информацию о пользователе", errorResponse["message"])
		assert.Equal(t, float64(http.StatusUnauthorized), errorResponse["statusCode"])

		mockAuthService.AssertExpectations(t)
	})

	t.Run("Me failed without claims", func(t *testing.T) {
		app := fiber.New()
		authHandler := NewAuthHandler(nil)
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("userClaims", nil)
			return c.Next()
		})
		app.Get("/auth/me", authHandler.Me)

		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var errorResponse map[string]interface{}
		assert.NoError(t, json.Unmarshal(body, &errorResponse))
		assert.Equal(t, "Не удалось получить информацию о пользователе", errorResponse["message"])
		assert.Equal(t, float64(http.StatusUnauthorized), errorResponse["statusCode"])
	})
}
