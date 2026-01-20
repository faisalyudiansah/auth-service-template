package converter

import (
	dto_response "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/response"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
)

func UserEntityToDTOLogin(user *entity.User) *dto_response.Login {
	if user == nil {
		return nil
	}

	return &dto_response.Login{
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsOauth:    user.IsOauth,
		IsActive:   user.IsActive,
		Role:       user.Role,
		RoleLabel:  user.Role.String(),
	}
}

func UserEntityToDTORegister(user *entity.User) *dto_response.Register {
	if user == nil {
		return nil
	}

	return &dto_response.Register{
		ID:         user.ID,
		Email:      user.Email,
		FullName:   user.UserDetail.FullName,
		IsVerified: user.IsVerified,
		IsOauth:    user.IsOauth,
		IsActive:   user.IsActive,
		Role:       user.Role,
		RoleLabel:  user.Role.String(),
		Sex:        user.UserDetail.Sex,
		SexLabel:   user.UserDetail.Sex.String(),
		BirthDate:  user.UserDetail.BirthDate,
		CreatedAt:  user.CreatedAt,
	}
}
