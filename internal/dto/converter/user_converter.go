package converter

import (
	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/dto"
)

func UserToResponse(user *ent.User) *dto.UserResponse {
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
