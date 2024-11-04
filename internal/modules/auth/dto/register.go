package dto

type RegisterDTO struct {
	Name     string  `json:"name" validate:"required,min=4"`
	NickName string  `json:"nickName" validate:"required,min=4"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Photo    *string `json:"photo" validate:"omitempty,url"`
}
