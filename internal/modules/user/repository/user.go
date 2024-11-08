package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/internal/modules/user/dto"
	"kbox-api/shared/httperror"
)

type UserRepositoryInterface interface {
	GetOneById(ctx context.Context, id string) (model.User, error)
	GetOneByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	GetOneByNickName(ctx context.Context, nickName string) (model.User, error)
	Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error)
}

type userRepository struct {
	storage *postgres.Storage
}

func NewUserRepository(storage *postgres.Storage) UserRepositoryInterface {
	return &userRepository{
		storage: storage,
	}
}

func (u *userRepository) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
	var user model.User
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
		return model.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *userRepository) GetOneById(ctx context.Context, id string) (model.User, error) {
	var user model.User

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
		return model.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	} else if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *userRepository) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	result := u.storage.DB().WithContext(ctx).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, nil
	} else if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}

	return user, nil
}

func (u *userRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	result := u.storage.DB().WithContext(ctx).Create(&user)
	if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *userRepository) Update(ctx context.Context, userDTO dto.UpdateUserDTO) (model.User, error) {
	tx := u.storage.DB().WithContext(ctx).Begin()

	var user model.User
	if err := tx.Where("id = ?", userDTO.ID).First(&user).Error; err != nil {
		tx.Rollback()
		return model.User{}, httperror.New(http.StatusNotFound, "Пользователь не найден")
	}

	if err := tx.Model(&user).Updates(userDTO).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key value") {
			return model.User{}, httperror.New(
				http.StatusConflict,
				"Пользователь с таким никнеймом или почтой уже существует",
			)
		}
		return model.User{}, httperror.New(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return user, nil
}
