package entity

import (
	custom_typeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	entityPkg "github.com/faisalyudiansah/auth-service-template/pkg/entity"
	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"

	"github.com/google/uuid"
)

type UserDetail struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	FullName    string
	Sex         custom_typeAuth.Sex
	PhoneNumber *string
	ImageURL    string
	BirthDate   custom_typePkg.DateOnly

	entityPkg.Audit
}
