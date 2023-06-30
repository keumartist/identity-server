package mocks

import (
	"art-sso/internal/service/token"

	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	token.TokenService
	mock.Mock
}

func (m *MockTokenService) GenerateToken(input token.GenerateTokenInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) VerifyToken(input token.VerifyTokenInput) (bool, string, string, error) {
	args := m.Called(input)
	return args.Bool(0), args.String(1), args.String(2), args.Error(3)
}
