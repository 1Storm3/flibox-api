package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"kinopoisk-api/shared/httperror"
	"net/http"
	"time"
)

type User struct {
	Id         string    `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	NickName   string    `json:"nickName" gorm:"column:nick_name"`
	Name       string    `json:"name" gorm:"column:name"`
	Email      string    `json:"email" gorm:"column:email"`
	Password   string    `json:"password" gorm:"column:password"`
	Photo      string    `json:"photo" gorm:"column:photo"`
	Role       string    `json:"role" gorm:"column:role"`
	IsVerified bool      `json:"isVerified" gorm:"column:is_verified"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetOneByNickName(nickName string) (User, error) {
	return s.userRepo.GetOneByNickName(context.Background(), nickName)
}

func (s *UserService) GetOneById(id string) (User, error) {
	return s.userRepo.GetOneById(context.Background(), id)
}

func (s *UserService) GetOneByEmail(email string) (User, error) {
	return s.userRepo.GetOneByEmail(context.Background(), email)
}

func (s *UserService) CheckPassword(user User, password string) bool {
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

func (s *UserService) CreateUser(user User) (User, error) {
	return s.userRepo.CreateUser(context.Background(), user)
}
