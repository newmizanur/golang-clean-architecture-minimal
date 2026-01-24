package converter

import (
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	"golang-clean-architecture/internal/model"
)

func UserToResponse(user *dbmodel.Users) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToTokenResponse(token string) *model.UserResponse {
	return &model.UserResponse{
		Token: token,
	}
}
