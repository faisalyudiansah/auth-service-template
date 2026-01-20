package middleware

import (
	"context"
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/utils"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/jwtutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/redisutils"
	sessioncookieutils "github.com/faisalyudiansah/auth-service-template/pkg/utils/sessionCookieUtils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	redisUtil redisutils.RedisUtil
	jwtUtil   jwtutils.JwtUtilInterface
}

func NewAuthMiddleware(
	redisUtil redisutils.RedisUtil,
	jwtUtil jwtutils.JwtUtilInterface,
) *AuthMiddleware {
	return &AuthMiddleware{
		redisUtil: redisUtil,
		jwtUtil:   jwtUtil,
	}
}

func (m *AuthMiddleware) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := sessioncookieutils.GetSessionIDFromCookie(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		if sessionID == uuid.Nil {
			ctx.Error(apperror.NewForbiddenAccessError())
			ctx.Abort()
			return
		}

		getSession := &entity.Session{}

		if err := m.redisUtil.GetWithScanJSON(
			ctx,
			utils.SessionKey(sessionID),
			getSession,
		); err != nil || getSession.UserID == uuid.Nil {
			ctx.Error(apperror.NewForbiddenAccessError())
			ctx.Abort()
			return
		}

		if getSession.AccessToken == "" || len(getSession.AccessToken) == 0 {
			ctx.Error(apperror.NewForbiddenAccessError())
			ctx.Abort()
			return
		}

		if ctx.FullPath() == "/auth/logout" || ctx.FullPath() == "/auth/refresh" {
			m.injectCtxSessionID(ctx, sessionID)
			ctx.Next()
			return
		}

		claims, err := m.jwtUtil.Parse(getSession.AccessToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.Error(apperror.NewExpiredTokenError())
				ctx.Abort()
				return
			} else {
				ctx.Error(apperror.NewForbiddenAccessError())
				ctx.Abort()
				return
			}
		}

		m.injectCtx(claims, ctx, sessionID)

		ctx.Next()
	}
}

func (m *AuthMiddleware) ProtectedRoles(allowedRoles ...custom_type.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := utils.GetValueRoleUserFromContext(ctx)

		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.Error(apperror.NewUnauthorizedError())
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *AuthMiddleware) OnlySelfOrAdmin(paramName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := utils.GetValueUserIDFromContext(ctx)
		role := utils.GetValueRoleUserFromContext(ctx)

		if role == custom_type.RoleAdmin {
			ctx.Next()
			return
		}

		paramUserID := ctx.Param(paramName)
		targetUserID, err := uuid.Parse(paramUserID)
		if err != nil || userID != targetUserID {
			ctx.Error(apperror.NewDontHavePermissionErrorMessageError())
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *AuthMiddleware) injectCtx(claims *jwtutils.JWTClaims, ctx *gin.Context, sessionID uuid.UUID) {
	m.injectCtxSessionID(ctx, sessionID)

	reqCtx := ctx.Request.Context()
	reqCtx = context.WithValue(reqCtx, constant.ContextUserID, claims.UserID)
	reqCtx = context.WithValue(reqCtx, constant.ContextRole, claims.Role)
	reqCtx = context.WithValue(reqCtx, constant.ContextJTI, claims.ID)

	ctx.Request = ctx.Request.WithContext(reqCtx)
}

func (m *AuthMiddleware) injectCtxSessionID(ctx *gin.Context, sessionID uuid.UUID) {
	reqCtx := context.WithValue(ctx.Request.Context(), constant.ContextSessionID, sessionID)
	ctx.Request = ctx.Request.WithContext(reqCtx)
}
