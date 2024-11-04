package service

import (
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/user/dto"
	"kbox-api/shared/httperror"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetOneByNickName(nickName string) (model.User, error) {
	return s.userRepo.GetOneByNickName(context.Background(), nickName)
}

func (s *UserService) GetOneById(id string) (model.User, error) {
	return s.userRepo.GetOneById(context.Background(), id)
}

func (s *UserService) GetOneByEmail(email string) (model.User, error) {
	return s.userRepo.GetOneByEmail(context.Background(), email)
}

func (s *UserService) CheckPassword(user model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return string(hashedPassword), nil
}

func (s *UserService) Create(user model.User) (model.User, error) {
	return s.userRepo.Create(context.Background(), user)
}

func (s *UserService) Update(userDTO dto.UpdateUserDTO) (model.User, error) {
	return s.userRepo.Update(context.Background(), userDTO)
}
