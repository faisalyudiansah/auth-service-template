package entity

import (
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"

	"github.com/google/uuid"
)

type Session struct {
	UserID       uuid.UUID        `json:"user_id"`
	Role         custom_type.Role `json:"role"`
	JTI          string           `json:"jti"`
	SessionID    uuid.UUID        `json:"session_id"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	LoginAt      uint64           `json:"login_at"`
}
