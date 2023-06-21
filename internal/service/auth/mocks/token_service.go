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
