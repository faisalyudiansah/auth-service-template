package converter

import (
	dto_response "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/response"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	converterPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto/converter"
)

func UserEntityToDTOResponse(e *entity.User) *dto_response.User {
	if e == nil {
		return nil
	}
	convert := &dto_response.User{
		ID:         e.ID,
		Role:       e.Role,
		RoleLabel:  e.Role.String(),
		Email:      e.Email,
		IsVerified: e.IsVerified,
		IsOauth:    e.IsOauth,
		IsActive:   e.IsActive,
		Audit:      converterPkg.ConvertAudit(e.Audit),
		UserDetail: UserDetailEntityToDTOResponse(e.UserDetail),
	}
	return convert
}

func UserDetailEntityToDTOResponse(e *entity.UserDetail) *dto_response.UserDetail {
	if e == nil {
		return nil
	}
	convert := &dto_response.UserDetail{
		ID:          e.ID,
		UserID:      e.UserID,
		FullName:    e.FullName,
		Sex:         e.Sex,
		SexLabel:    e.Sex.String(),
		PhoneNumber: e.PhoneNumber,
		ImageURL:    e.ImageURL,
		BirthDate:   e.BirthDate,
		Audit:       converterPkg.ConvertAudit(e.Audit),
	}
	return convert
}

func ListUserEntityToDTOResponse(e []*entity.User) []dto_response.User {
	result := make([]dto_response.User, 0, len(e))

	for _, item := range e {
		if item == nil {
			continue
		}
		dto := UserEntityToDTOResponse(item)
		if dto != nil {
			result = append(result, *dto)
		}
	}

	return result
}
