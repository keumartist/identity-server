package mocks

import (
	"art-sso/internal/service/mail"

	"github.com/stretchr/testify/mock"
)

type MockMailService struct {
	mail.MailService
	mock.Mock
}

func (m *MockMailService) SendVerificationEmail(email, verificationCode string) error {
	args := m.Called(email, verificationCode)
	return args.Error(0)
}
