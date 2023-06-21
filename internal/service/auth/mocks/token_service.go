package mocks

import (
	"art-sso/internal/service/token"

	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	token.TokenService
	mock.Mock
}

func (m *MockTokenService) GenerateIdToken(userID, email string, expirationInSeconds uint) (string, error) {
	args := m.Called(userID, email, expirationInSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID, email string, expirationInSeconds uint) (string, error) {
	args := m.Called(userID, email, expirationInSeconds)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateAccessToken(userID, email string, expirationInSeconds uint) (string, error) {
	args := m.Called(userID, email, expirationInSeconds)
	return args.String(0), args.Error(1)
}
