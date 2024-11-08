package service

import (
	"context"

	"net/http"

	"golang.org/x/crypto/bcrypt"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/external/service"
	"kbox-api/internal/modules/user/dto"
	"kbox-api/internal/modules/user/repository"
	"kbox-api/shared/helper"
	"kbox-api/shared/httperror"
)

type UserServiceInterface interface {
	GetOneByNickName(nickName string) (model.User, error)
	GetOneByEmail(email string) (model.User, error)
	CheckPassword(user model.User, password string) bool
	HashPassword(password string) (string, error)
	Create(user model.User) (model.User, error)
	GetOneById(id string) (model.User, error)
	Update(userDTO dto.UpdateUserDTO) (model.User, error)
}

type UserService struct {
	userRepo  repository.UserRepositoryInterface
	s3Service service.S3ServiceInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface, s3Service service.S3ServiceInterface) UserServiceInterface {
	return &UserService{
		userRepo:  userRepo,
		s3Service: s3Service,
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
	if userDTO.Photo != nil {
		user, err := s.GetOneById(userDTO.ID)
		if err != nil {
			return model.User{}, err
		}

		if user.Photo != nil && *user.Photo != "" {
			key, err := helper.ExtractS3Key(*user.Photo)
			if err != nil {
				return model.User{}, err
			}

			err = s.s3Service.DeleteFile(context.Background(), key)
			if err != nil {
				return model.User{}, err
			}
		}
	}

	return s.userRepo.Update(context.Background(), userDTO)
}
