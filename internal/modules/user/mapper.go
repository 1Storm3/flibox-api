package user

import "github.com/1Storm3/flibox-api/internal/model"

func MapModelUserToResponseDTO(user model.User) ResponseDTO {
	return ResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		NickName:   user.NickName,
		Email:      user.Email,
		Photo:      user.Photo,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		IsVerified: user.IsVerified,
	}
}
