package auth

import (
	"kbox-api/internal/model"
)

func MapModelUserToResponseDTO(user model.User) MeResponseDTO {
	return MeResponseDTO{
		Id:         user.ID,
		Name:       user.Name,
		NickName:   user.NickName,
		Email:      user.Email,
		Photo:      user.Photo,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt.String(),
		UpdatedAt:  user.UpdatedAt.String(),
		IsVerified: user.IsVerified,
	}
}
