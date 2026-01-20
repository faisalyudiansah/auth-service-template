package entity

import (
	"time"

	"github.com/google/uuid"
)

type Audit struct {
	CreatedAt     time.Time
	CreatedBy     uuid.UUID
	UpdatedAt     *time.Time
	UpdatedBy     *uuid.UUID
	DeletedAt     *time.Time
	DeletedBy     *uuid.UUID
	DeletedReason *string
}
