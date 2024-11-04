package mapper

import (
	"kbox-api/internal/model"
	"kbox-api/internal/modules/auth/dto"
)

func MapModelUserToResponseDTO(user model.User) dto.MeResponseDTO {
	return dto.MeResponseDTO{
		Id:         user.Id,
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
