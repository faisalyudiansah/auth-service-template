package dto

import (
	"time"

	"github.com/google/uuid"
)

type Audit struct {
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy uuid.UUID  `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy *uuid.UUID `json:"deleted_by"`
}
