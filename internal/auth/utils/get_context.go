package utils

import (
	"context"

	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/pkg/constant"

	"github.com/google/uuid"
)

func GetValueSessionIDFromContext(c context.Context) uuid.UUID {
	if userId, ok := c.Value(constant.ContextSessionID).(uuid.UUID); ok {
		return userId
	}
	return uuid.Nil
}

func GetValueUserIDFromContext(c context.Context) uuid.UUID {
	if userId, ok := c.Value(constant.ContextUserID).(uuid.UUID); ok {
		return userId
	}
	return uuid.Nil
}

func GetValueRoleUserFromContext(c context.Context) custom_type.Role {
	if role, ok := c.Value(constant.ContextRole).(custom_type.Role); ok {
		return role
	}
	return 0
}

func GetJTIFromContext(c context.Context) string {
	if jti, ok := c.Value(constant.ContextJTI).(string); ok {
		return jti
	}
	return ""
}
