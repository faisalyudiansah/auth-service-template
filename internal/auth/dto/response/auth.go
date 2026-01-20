package dto_response

import (
	"time"

	custom_typeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"

	"github.com/google/uuid"
)

type Login struct {
	Email      string               `json:"email"`
	IsVerified bool                 `json:"is_verified"`
	IsOauth    bool                 `json:"is_outh"`
	IsActive   bool                 `json:"is_active"`
	Role       custom_typeAuth.Role `json:"role"`
	RoleLabel  string               `json:"role_label"`
}

type Register struct {
	ID         uuid.UUID               `json:"id"`
	Email      string                  `json:"email"`
	FullName   string                  `json:"full_name"`
	IsVerified bool                    `json:"is_verified"`
	IsOauth    bool                    `json:"is_oauth"`
	IsActive   bool                    `json:"is_active"`
	Role       custom_typeAuth.Role    `json:"role"`
	RoleLabel  string                  `json:"role_label"`
	Sex        custom_typeAuth.Sex     `json:"sex"`
	SexLabel   string                  `json:"sex_label"`
	BirthDate  custom_typePkg.DateOnly `json:"birth_date"`
	CreatedAt  time.Time               `json:"created_at"`
}
