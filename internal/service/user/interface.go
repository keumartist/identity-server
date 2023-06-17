package user

import (
	dto "art-sso/internal/dto/user"
)

type UserService interface {
	CreateUser(input CreateUserInput) (dto.User, error)
	GetUserByID(input GetUserByIDInput) (dto.User, error)
	GetUserByEmail(input GetUserByEmailInput) (dto.User, error)
	UpdateUser(input UpdateUserInput) error
	DeleteUser(input DeleteUserInput) error
}
