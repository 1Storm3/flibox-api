package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/modules/user/service"
	"kinopoisk-api/shared/httperror"
	"net/http"
)

type userRepository struct {
	storage *postgres.Storage
}

func NewUserRepository(storage *postgres.Storage) *userRepository {
	return &userRepository{
		storage: storage,
	}
}

func (u *userRepository) GetOneByNickName(ctx context.Context, nickName string) (service.User, error) {
	var user service.User
	result := u.storage.DB().WithContext(ctx).
		Select("id",
			"nick_name",
			"name",
			"email",
			"photo",
			"role",
			"is_verified",
			"updated_at",
			"created_at",
		).
		Where("nick_name = ?", nickName).
		First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return service.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return service.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *userRepository) GetOneById(ctx context.Context, id string) (service.User, error) {
	var user service.User

	result := u.storage.DB().WithContext(ctx).
		Select("id",
			"nick_name",
			"name",
			"email",
			"photo",
			"role",
			"is_verified",
			"updated_at",
			"created_at",
		).
		Where("id = ?", id).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return service.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return service.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *userRepository) GetOneByEmail(ctx context.Context, email string) (service.User, error) {
	var user service.User

	result := u.storage.DB().WithContext(ctx).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return service.User{}, nil
	} else if result.Error != nil {
		return service.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user service.User) (service.User, error) {
	result := u.storage.DB().WithContext(ctx).Create(&user)
	if result.Error != nil {
		return service.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}
