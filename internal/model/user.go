package model

import "time"

type User struct {
	Id         string    `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	NickName   string    `json:"nickName" gorm:"column:nick_name"`
	Name       string    `json:"name" gorm:"column:name"`
	Email      string    `json:"email" gorm:"column:email"`
	Password   string    `json:"password" gorm:"column:password"`
	Photo      *string   `json:"photo" gorm:"column:photo"`
	Role       string    `json:"role" gorm:"column:role"`
	IsVerified bool      `json:"isVerified" gorm:"column:is_verified"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
}
