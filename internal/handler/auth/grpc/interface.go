package grpc

import (
	api "art-sso/internal/grpc/api/proto"
	"context"
)

type AuthHandler interface {
	VerifyToken(ctx context.Context, req *api.VerifyTokenRequest) (*api.VerifyTokenResponse, error)
}
