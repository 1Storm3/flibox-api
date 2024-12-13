package user

import (
	"context"

	"net/http"

	"golang.org/x/crypto/bcrypt"

	"kbox-api/internal/model"
	"kbox-api/internal/modules/external"
	"kbox-api/internal/shared/helper"
	"kbox-api/internal/shared/httperror"
)

var _ ServiceInterface = (*Service)(nil)

type ServiceInterface interface {
	GetOneByNickName(ctx context.Context, nickName string) (model.User, error)
	GetOneByEmail(ctx context.Context, email string) (model.User, error)
	CheckPassword(ctx context.Context, user model.User, password string) bool
	HashPassword(ctx context.Context, password string) (string, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	GetOneById(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, userDTO UpdateUserDTO) (model.User, error)
	UpdateForVerify(ctx context.Context, userDTO UpdateForVerifyDTO) (model.User, error)
}

type Service struct {
	repository RepositoryInterface
	s3Service  external.S3ServiceInterface
}

func NewUserService(repository RepositoryInterface, s3Service external.S3ServiceInterface) *Service {
	return &Service{
		repository: repository,
		s3Service:  s3Service,
	}
}

func (s *Service) UpdateForVerify(ctx context.Context, userDTO UpdateForVerifyDTO) (model.User, error) {
	return s.repository.UpdateForVerify(ctx, userDTO)
}

func (s *Service) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
	return s.repository.GetOneByNickName(ctx, nickName)
}

func (s *Service) GetOneById(ctx context.Context, id string) (model.User, error) {

	user, err := s.repository.GetOneById(ctx, id)

	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *Service) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
	return s.repository.GetOneByEmail(ctx, email)
}

func (s *Service) CheckPassword(_ context.Context, user model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *Service) HashPassword(_ context.Context, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", httperror.New(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	return string(hashedPassword), nil
}

func (s *Service) Create(ctx context.Context, user model.User) (model.User, error) {
	return s.repository.Create(ctx, user)
}

func (s *Service) Update(ctx context.Context, userDTO UpdateUserDTO) (model.User, error) {
	if userDTO.Photo != nil {
		user, err := s.GetOneById(ctx, userDTO.ID)
		if err != nil {
			return model.User{}, err
		}

		if user.Photo != nil && *user.Photo != "" {
			key, err := helper.ExtractS3Key(*user.Photo)
			if err != nil {
				return model.User{}, err
			}

			err = s.s3Service.DeleteFile(ctx, key)
			if err != nil {
				return model.User{}, err
			}
		}
	}

	return s.repository.Update(context.Background(), userDTO)
}
