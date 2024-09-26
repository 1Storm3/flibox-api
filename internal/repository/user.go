package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/service"
)

type userRepository struct {
	storage *postgres.Storage
}

func NewUserRepository(storage *postgres.Storage) *userRepository {
	return &userRepository{
		storage: storage,
	}
}

func (u *userRepository) GetOne(ctx context.Context, userToken string) (service.User, error) {
	var user service.User

	result := u.storage.DB().WithContext(ctx).Where("user_token = ?", userToken).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return service.User{}, nil
	} else if result.Error != nil {
		return service.User{}, result.Error
	}

	return user, nil
}
