package service

import "context"

type User struct {
	Uuid      string `json:"uuid" gorm:"column:uuid"`
	FirstName string `json:"firstName" gorm:"column:first_name"`
	Avatar    string `json:"avatar" gorm:"column:avatar"`
	TgId      string `json:"tgId" gorm:"column:tg_id"`
	UserToken string `json:"userToken" gorm:"column:user_token"`
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetOne(userToken string) (User, error) {
	return s.userRepo.GetOne(context.Background(), userToken)
}
