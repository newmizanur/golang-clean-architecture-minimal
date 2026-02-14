package converter

import (
	"golang-clean-architecture/internal/dto"
	dbmodel "golang-clean-architecture/internal/persistence/model"
)

func UserToResponse(user *dbmodel.Users) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToTokenResponse(token string) *dto.UserResponse {
	return &dto.UserResponse{
		Token: token,
	}
}
