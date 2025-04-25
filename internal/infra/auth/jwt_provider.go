package auth

import (
	"backend/internal/domain/auth"
	"context"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

var _ auth.TokenProvider[auth.Claims] = (*JwtProvider)(nil)

type JwtProvider struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJwtProvider(publicKey, privateKey string) (*JwtProvider, error) {
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "parse public key error")
	}
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "parse private key error")
	}
	return &JwtProvider{
		PrivateKey: priKey,
		PublicKey:  pubKey,
	}, nil
}

func (j *JwtProvider) Generate(_ context.Context, claims *auth.Claims) (token string, err error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)
	token, err = t.SignedString(j.PrivateKey)
	return
}

func (j *JwtProvider) Verify(_ context.Context, tokenString string) (claims *auth.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signature signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.PublicKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "invalid token")
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.Wrap(err, "invalid token")
}
