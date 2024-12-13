package user

import "time"

type CreateUserDTO struct {
	Name     string  `json:"name" validate:"required,min=4"`
	NickName string  `json:"nickName" validate:"required,min=4"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Photo    *string `json:"photo" validate:"omitempty,url"`
}

type ResponseDTO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NickName   string    `json:"nickName"`
	Email      string    `json:"email"`
	Photo      *string   `json:"photo"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	IsVerified bool      `json:"isVerified"`
}

type UpdateUserDTO struct {
	ID           string  `json:"id" validate:"required"`
	NickName     *string `json:"nickName" validate:"omitempty,min=4"`
	Name         *string `json:"name" validate:"omitempty,min=4"`
	Email        *string `json:"email" validate:"omitempty,email"`
	LastActivity *string `json:"lastActivity" validate:"omitempty"`
	Photo        *string `json:"photo" validate:"omitempty,url"`
}

type UpdateForVerifyDTO struct {
	ID            string  `json:"id" validate:"required"`
	IsVerified    bool    `json:"isVerified" validate:"required"`
	VerifiedToken *string `json:"verifiedToken" validate:"required"`
}
