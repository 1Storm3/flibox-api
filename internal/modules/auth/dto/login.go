package dto

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
