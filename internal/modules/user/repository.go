package user

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"kbox-api/database/postgres"
	"kbox-api/internal/model"
	"kbox-api/internal/shared/httperror"
)

var _ RepositoryInterface = (*Repository)(nil)

type RepositoryInterface interface {
	GetOneById(ctx context.Context, id string) (model.User, error)
	GetOneByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	GetOneByNickName(ctx context.Context, nickName string) (model.User, error)
	Update(ctx context.Context, userDTO UpdateUserDTO) (model.User, error)
	UpdateForVerify(ctx context.Context, userDTO UpdateForVerifyDTO) (model.User, error)
}

type Repository struct {
	storage *postgres.Storage
}

func NewUserRepository(storage *postgres.Storage) RepositoryInterface {
	return &Repository{
		storage: storage,
	}
}

func (u *Repository) UpdateForVerify(ctx context.Context, userDTO UpdateForVerifyDTO) (model.User, error) {
	var user model.User
	result := u.storage.DB().WithContext(ctx).Where("id = ?", userDTO.ID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{},
			httperror.New(
				http.StatusNotFound,
				"Пользователь не найден",
			)
	}

	user.IsVerified = userDTO.IsVerified
	user.VerifiedToken = userDTO.VerifiedToken

	result = u.storage.DB().WithContext(ctx).Save(&user)
	if result.Error != nil {
		return model.User{},
			httperror.New(
				http.StatusInternalServerError,
				result.Error.Error(),
			)
	}
	return user, nil
}

func (u *Repository) GetOneByNickName(ctx context.Context, nickName string) (model.User, error) {
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

func (u *Repository) GetOneById(ctx context.Context, id string) (model.User, error) {
	var user model.User

	result := u.storage.DB().
		Select("id",
			"nick_name",
			"name",
			"email",
			"photo",
			"role",
			"is_verified",
			"is_blocked",
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

func (u *Repository) GetOneByEmail(ctx context.Context, email string) (model.User, error) {
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

func (u *Repository) Create(ctx context.Context, user model.User) (model.User, error) {
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

func (u *Repository) Update(ctx context.Context, userDTO UpdateUserDTO) (model.User, error) {
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
