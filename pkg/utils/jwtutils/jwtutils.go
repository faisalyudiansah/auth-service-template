package jwtutils

import (
	"errors"
	"time"

	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtUtilInterface interface {
	Sign(userID uuid.UUID, role custom_type.Role, jti string, currentTime time.Time) (string, error)
	SignRefresh(currentTime time.Time) (string, error)
	Parse(tokenString string) (*JWTClaims, error)
}

type jwtUtil struct {
	jwtConfig *config.JwtConfig
}

func NewJwtUtil(jwtConfig *config.JwtConfig) *jwtUtil {
	return &jwtUtil{
		jwtConfig: jwtConfig,
	}
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID  uuid.UUID        `json:"user_id"`
	Role    custom_type.Role `json:"role"`
	LoginAt uint64           `json:"login_at"`
}

func (h *jwtUtil) Sign(userID uuid.UUID, role custom_type.Role, jti string, currentTime time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID:  userID,
		Role:    role,
		LoginAt: uint64(currentTime.UnixMilli()),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(h.jwtConfig.TokenDuration) * time.Minute)),
			Issuer:    h.jwtConfig.Issuer,
		},
	})

	s, err := token.SignedString([]byte(h.jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return s, nil
}

func (h *jwtUtil) SignRefresh(currentTime time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(h.jwtConfig.RefreshDuration) * time.Minute)),
			Issuer:    h.jwtConfig.Issuer,
		},
	})

	s, err := token.SignedString([]byte(h.jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return s, nil
}

func (h *jwtUtil) Parse(tokenString string) (*JWTClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods(h.jwtConfig.AllowedAlgs),
		jwt.WithIssuer(h.jwtConfig.Issuer),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	return h.parseClaims(parser, tokenString)
}

func (h *jwtUtil) parseClaims(parser *jwt.Parser, tokenString string) (*JWTClaims, error) {
	token, err := parser.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(h.jwtConfig.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token not valid")
	}

	return claims, nil
}
