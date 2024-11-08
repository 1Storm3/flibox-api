package mapper

import (
	"kbox-api/internal/model"
	"kbox-api/internal/modules/user/dto"
)

func MapModelUserToResponseDTO(user model.User) dto.CreateUserResponseDTO {
	return dto.CreateUserResponseDTO{
		ID:         user.ID,
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
