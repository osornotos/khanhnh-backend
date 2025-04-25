package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"khanhnh-backend/internal/domain/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"khanhnh-backend/tests/mocks"
)

func TestUserServiceImpl_Login(t *testing.T) {
	mockTokenProvider := new(mocks.TokenProviderMock[auth.Claims])
	service := NewUserServiceImpl(mockTokenProvider)

	ctx := context.Background()
	expectedToken := "mocked-token"

	// Test successful token generation
	mockTokenProvider.On("Generate", ctx, mock.AnythingOfType("*auth.Claims")).Return(expectedToken, nil)

	token, err := service.Login(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockTokenProvider.AssertExpectations(t)
}
