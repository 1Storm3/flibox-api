package service

import "context"

type User struct {
	Id        string `json:"id" gorm:"column:id"`
	FirstName string `json:"firstName" gorm:"column:first_name"`
	Avatar    string `json:"avatar" gorm:"column:avatar"`
	TgId      string `json:"tgId" gorm:"column:tg_id"`
	UserToken string `json:"userToken" gorm:"column:user_token"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at"`
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
