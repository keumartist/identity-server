package mocks

import (
	dto "art-sso/internal/dto/user"
	userservice "art-sso/internal/service/user"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(input userservice.CreateUserInput) (dto.User, error) {
	args := m.Called(input)
	return args.Get(0).(dto.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(input userservice.GetUserByIDInput) (dto.User, error) {
	args := m.Called(input)
	return args.Get(0).(dto.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(input userservice.GetUserByEmailInput) (dto.User, error) {
	args := m.Called(input)
	return args.Get(0).(dto.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(input userservice.UpdateUserInput) error {
	return m.Called(input).Error(0)
}

func (m *MockUserService) DeleteUser(input userservice.DeleteUserInput) error {
	return m.Called(input).Error(0)
}
