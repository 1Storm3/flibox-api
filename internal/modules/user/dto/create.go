package dto

type CreateUserDTO struct {
	Name     string  `json:"name" validate:"required,min=4"`
	NickName string  `json:"nickName" validate:"required,min=4"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Photo    *string `json:"photo" validate:"omitempty,url"`
}

type CreateUserResponseDTO struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	NickName   string  `json:"nickName"`
	Email      string  `json:"email"`
	Photo      *string `json:"photo"`
	Role       string  `json:"role"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	IsVerified bool    `json:"isVerified"`
}
