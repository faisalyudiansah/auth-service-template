package entity

import (
	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"

	"github.com/google/uuid"
)

type VerificationToken struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	VerificationToken uuid.UUID

	dtoPkg.Audit
}
