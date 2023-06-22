package grpc_test

import (
	"context"
	"testing"

	api "art-sso/internal/grpc/api/proto"
	authhandler "art-sso/internal/handler/auth/grpc"
	mocks "art-sso/internal/handler/auth/grpc/mocks"
	"art-sso/internal/service/token"

	"github.com/stretchr/testify/assert"
)

func setup() (*mocks.MockTokenService, authhandler.AuthHandler) {
	mockTokenService := new(mocks.MockTokenService)
	handler := authhandler.NewAuthHandler(mockTokenService)

	return mockTokenService, handler
}

func TestVerifyToken(t *testing.T) {
	t.Run("Valid token", func(t *testing.T) {
		mockTokenService, handler := setup()

		mockTokenService.On("VerifyToken", token.VerifyTokenInput{
			Token:     "validToken",
			TokenType: token.AccessToken,
		}).Return(true, "userId", "email@example.com", nil)

		req := &api.VerifyTokenRequest{
			Token:     "validToken",
			TokenType: "access",
		}

		resp, err := handler.VerifyToken(context.Background(), req)

		assert.NoError(t, err)
		assert.True(t, resp.IsValid)
		assert.Equal(t, "userId", resp.UserId)
		assert.Equal(t, "email@example.com", resp.Email)
	})

	t.Run("Invalid token", func(t *testing.T) {
		mockTokenService, handler := setup()

		mockTokenService.On("VerifyToken", token.VerifyTokenInput{
			Token:     "invalidToken",
			TokenType: token.AccessToken,
		}).Return(false, "", "", nil)

		req := &api.VerifyTokenRequest{
			Token:     "invalidToken",
			TokenType: "access",
		}

		resp, err := handler.VerifyToken(context.Background(), req)

		assert.NoError(t, err)
		assert.False(t, resp.IsValid)
		assert.Equal(t, "", resp.UserId)
		assert.Equal(t, "", resp.Email)
	})
}
