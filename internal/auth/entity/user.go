package entity

import (
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	entityPkg "github.com/faisalyudiansah/auth-service-template/pkg/entity"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Role         custom_type.Role
	Email        string
	HashPassword string
	IsVerified   bool
	IsOauth      bool
	IsActive     bool

	entityPkg.Audit

	UserDetail *UserDetail
}
