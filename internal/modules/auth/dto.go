package auth

// LoginDTO dtoAuth.LoginDTO представляет данные для входа пользователя
// @swagger:model
type LoginDTO struct {
	// Email пользователя
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`

	// Пароль пользователя
	// required: true
	// min length: 6
	// example: password123
	Password string `json:"password" validate:"required,min=6"`
}

type MeResponseDTO struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	NickName   string  `json:"nickName"`
	Email      string  `json:"email"`
	Photo      *string `json:"photo"`
	Role       string  `json:"role"`
	IsVerified bool    `json:"isVerified"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type RegisterDTO struct {
	Name     string  `json:"name" validate:"required,min=4"`
	NickName string  `json:"nickName" validate:"required,min=4"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Photo    *string `json:"photo" validate:"omitempty,url"`
}
