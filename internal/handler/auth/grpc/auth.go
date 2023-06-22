package grpc

import (
	"context"
	"fmt"

	api "art-sso/internal/grpc/api/proto"

	"art-sso/internal/service/token"
	service "art-sso/internal/service/token"
)

type AuthHandlerImpl struct {
	tokenService service.TokenService
}

func NewAuthHandler(tokenService service.TokenService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		tokenService: tokenService,
	}
}

func (h *AuthHandlerImpl) VerifyToken(ctx context.Context, req *api.VerifyTokenRequest) (*api.VerifyTokenResponse, error) {
	var tokenType token.TokenType

	switch req.TokenType {
	case "access":
		tokenType = token.AccessToken
	case "id":
		tokenType = token.IdToken
	case "refresh":
		tokenType = token.RefreshToken
	default:
		return nil, fmt.Errorf("unknown token type: %s", req.TokenType)
	}

	input := token.VerifyTokenInput{
		Token:     req.Token,
		TokenType: tokenType,
	}

	isValid, userId, email, err := h.tokenService.VerifyToken(input)

	if err != nil {
		return nil, err
	}

	if !isValid {
		return &api.VerifyTokenResponse{
			IsValid: false,
			UserId:  "",
			Email:   "",
		}, nil
	}

	return &api.VerifyTokenResponse{
		IsValid: true,
		UserId:  userId,
		Email:   email,
	}, nil
}
