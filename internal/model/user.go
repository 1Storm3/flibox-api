package model

import "time"

type User struct {
	ID            string    `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	NickName      string    `json:"nickName" gorm:"column:nick_name"`
	Name          string    `json:"name" gorm:"column:name"`
	Email         string    `json:"email" gorm:"column:email"`
	Password      string    `json:"password" gorm:"column:password"`
	Photo         *string   `json:"photo" gorm:"column:photo"`
	Role          string    `json:"role" gorm:"column:role"`
	VerifiedToken *string   `json:"verifiedToken" gorm:"column:verified_token"`
	IsVerified    bool      `json:"isVerified" gorm:"column:is_verified"`
	IsBlocked     bool      `json:"isBlocked" gorm:"column:is_blocked"`
	LastActivity  time.Time `json:"lastActivity" gorm:"column:last_activity"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
}
