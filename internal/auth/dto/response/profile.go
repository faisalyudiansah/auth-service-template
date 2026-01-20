package dto_response

import (
	custom_typeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"
	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID            `json:"id"`
	Role       custom_typeAuth.Role `json:"role"`
	RoleLabel  string               `json:"role_label"`
	Email      string               `json:"email"`
	IsVerified bool                 `json:"is_verified"`
	IsOauth    bool                 `json:"is_oauth"`
	IsActive   bool                 `json:"is_active"`

	dtoPkg.Audit

	UserDetail *UserDetail `json:"user_detail"`
}

type UserDetail struct {
	ID          uuid.UUID               `json:"id"`
	UserID      uuid.UUID               `json:"user_id"`
	FullName    string                  `json:"full_name"`
	Sex         custom_typeAuth.Sex     `json:"sex"`
	SexLabel    string                  `json:"sex_label"`
	PhoneNumber *string                 `json:"phone_number"`
	ImageURL    string                  `json:"image_url"`
	BirthDate   custom_typePkg.DateOnly `json:"birth_date"`

	dtoPkg.Audit
}
