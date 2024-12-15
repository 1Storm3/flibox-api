package auth

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/1Storm3/flibox-api/internal/config"
	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/modules/user"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
	args := m.Called(ctx, nickName)
	return *args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserService) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return *args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserService) CheckPassword(ctx context.Context, user *model.User, password string) bool {
	args := m.Called(ctx, user, password)
	return args.Bool(0)
}

func (m *mockUserService) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *mockUserService) Create(ctx context.Context, user model.User) (model.User, error) {
	args := m.Called(ctx, user)
	return *args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserService) GetOneById(ctx context.Context, id string) (model.User, error) {
	args := m.Called(ctx, id)
	return *args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserService) Update(ctx context.Context, userDTO user.UpdateUserDTO) (model.User, error) {
	args := m.Called(ctx, userDTO)
	return *args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserService) UpdateForVerify(ctx context.Context, userDTO user.UpdateForVerifyDTO) (model.User, error) {
	args := m.Called(ctx, userDTO)
	return *args.Get(0).(*model.User), args.Error(1)
}

func TestLoginService(t *testing.T) {
	ctx := context.Background()
	mockUserService := new(mockUserService)

	s := Service{
		userService: mockUserService,
		cfg: &config.Config{
			App: config.AppConfig{
				JwtSecretKey: "secret",
				JwtExpiresIn: "1h",
			},
		},
	}

	tests := []struct {
		name          string
		mockUser      *model.User
		mockEmail     string
		mockPassword  string
		expectedError bool
		expectedCode  int
		expectedToken bool
		mockError     error
		checkPassword bool
	}{
		{
			name:          "success",
			mockUser:      &model.User{ID: "1", Role: "user", Email: "test@example.com"},
			mockEmail:     "test@example.com",
			mockPassword:  "password",
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedToken: true,
			mockError:     nil,
			checkPassword: true,
		},
		{
			name:          "invalid email or password",
			mockUser:      nil,
			mockEmail:     "test@example.com",
			mockPassword:  "password",
			expectedError: true,
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
			mockError:     httperror.New(http.StatusUnauthorized, "Неверный логин или пароль"),
			checkPassword: false,
		},
		{
			name:          "invalid password",
			mockUser:      &model.User{ID: "1", Role: "user", Email: "test@example.com"},
			mockEmail:     "test@example.com",
			mockPassword:  "wrong-password",
			expectedError: true,
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
			mockError:     nil,
			checkPassword: false,
		},
		{
			name:          "jwt generation failure",
			mockUser:      &model.User{ID: "1", Role: "user", Email: "test@example.com"},
			mockEmail:     "test@example.com",
			mockPassword:  "password",
			expectedError: true,
			expectedCode:  http.StatusInternalServerError,
			expectedToken: false,
			mockError:     nil,
			checkPassword: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockUser != nil {
				mockUserService.On("GetOneByEmail", ctx, tt.mockEmail).Return(tt.mockUser, tt.mockError)
				mockUserService.On("CheckPassword", ctx, tt.mockUser, tt.mockPassword).Return(tt.checkPassword)
			} else {
				mockUserService.On("GetOneByEmail", ctx, tt.mockEmail).Return(nil, tt.mockError)
			}

			if tt.name == "jwt generation failure" {
				s.cfg.App.JwtSecretKey = ""
			} else {
				s.cfg.App.JwtSecretKey = "secret"
			}
			result, err := s.Login(ctx, LoginDTO{Email: tt.mockEmail, Password: tt.mockPassword})

			if tt.expectedError {
				assert.Error(t, err)
				var httpErr *httperror.Error
				if errors.As(err, &httpErr) {
					assert.Equal(t, tt.expectedCode, httpErr.Code)
				}
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}

			mockUserService.AssertExpectations(t)
		})
	}
}
