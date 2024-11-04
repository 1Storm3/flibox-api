package dto

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
