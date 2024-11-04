package dto

type UpdateUserDTO struct {
	Id       string  `json:"id" validate:"required"`
	NickName *string `json:"nickName" validate:"omitempty,min=4"`
	Name     *string `json:"name" validate:"omitempty,min=4"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Photo    *string `json:"photo" validate:"omitempty,url"`
}
