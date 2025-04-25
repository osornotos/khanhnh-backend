package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider[Claims any] interface {
	Generate(ctx context.Context, claims *Claims) (token string, err error)
	Verify(ctx context.Context, token string) (claims *Claims, err error)
}

// Guest account

const (
	DefaultExpTime = 8 * time.Hour
	AdminRole      = "admin"
	UserRole       = "user"
)

type Claims struct {
	jwt.RegisteredClaims

	Roles []string `json:"roles"`
	Id    int64    `json:"id"`
}

func NewClaims(roles []string) Claims {
	now := time.Now().UTC()
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(DefaultExpTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "",
		},
		Roles: roles,
	}
}
