package dto_request

import (
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"

	"github.com/google/uuid"
)

type UpdateUser struct {
	UserID     uuid.UUID        `json:"user_id"`
	Role       custom_type.Role `json:"role" binding:"oneof=0 1 2"`
	IsVerified bool             `json:"is_verified" binding:"boolean"`
	IsActive   bool             `json:"is_active" binding:"boolean"`

	FullName    string          `json:"full_name" binding:"required"`
	Sex         custom_type.Sex `json:"sex" binding:"required,oneof=0 1 2"`
	PhoneNumber string          `json:"phone_number" binding:"required,phone_number"`
	ImageURL    string          `json:"image_url" binding:"required,url"`

	UpdatedBy     uuid.UUID        `json:"-"`
	RoleWhoIsEdit custom_type.Role `json:"-"`
}

type DeleteUser struct {
	UserID        uuid.UUID `json:"user_id"`
	DeletedReason string    `json:"deleted_reason" binding:"required"`

	DeletedBy uuid.UUID `json:"-"`
}
