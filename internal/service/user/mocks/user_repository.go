package mocks

import (
	domain "art-sso/internal/domain/user"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUnverifiedUser(user *domain.User, verificationCode string, expireAt time.Time) error {
	args := m.Called(user, verificationCode, expireAt)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserProfile(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRefreshToken(user *domain.User, refreshToken string) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateVerificationCode(user *domain.User, verificationCode string, expireAt time.Time) error {
	args := m.Called(user, verificationCode, expireAt)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) VerifyUserEmail(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
