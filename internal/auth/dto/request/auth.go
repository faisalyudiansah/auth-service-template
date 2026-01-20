package dto_request

import (
	custom_typeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"

	"github.com/google/uuid"
)

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Email     string                  `json:"email" binding:"required,email"`
	Password  string                  `json:"password" binding:"required,password"`
	FullName  string                  `json:"full_name" binding:"required"`
	Sex       custom_typeAuth.Sex     `json:"sex" binding:"required,oneof=0 1 2"`
	Role      custom_typeAuth.Role    `json:"role" binding:"oneof=0 1 2"`
	BirthDate custom_typePkg.DateOnly `json:"birth_date" binding:"required,birth_date"`

	CreatedBy *uuid.UUID `json:"-"`
}

type SendVerification struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyAccount struct {
	VerificationToken string `json:"verification_token" binding:"required"`
	Email             string `json:"email" binding:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPassword struct {
	ResetToken string `json:"reset_token" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required,password"`
}

type InactiveAccount struct {
	UserID    uuid.UUID `json:"user_id" binding:"required,user_id"`
	UpdatedBy uuid.UUID `json:"-"`
}
