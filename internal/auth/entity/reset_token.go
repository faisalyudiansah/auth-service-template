package entity

import (
	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"

	"github.com/google/uuid"
)

type ResetToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	ResetToken uuid.UUID

	dtoPkg.Audit
}
