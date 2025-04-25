package service

import (
	"context"
	"khanhnh-backend/internal/domain/auth"
)

type UserServiceImpl struct {
	tokenProvider auth.TokenProvider[auth.Claims]
}

func NewUserServiceImpl(
	tokenProvider auth.TokenProvider[auth.Claims],
) *UserServiceImpl {
	return &UserServiceImpl{
		tokenProvider: tokenProvider,
	}
}

func (u *UserServiceImpl) Login(ctx context.Context) (string, error) {
	// FIXME validate user info here
	claims := auth.NewClaims([]string{auth.UserRole})
	token, err := u.tokenProvider.Generate(ctx, &claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
