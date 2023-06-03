package repository

import domain "art-sso/internal/domain/user"

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(user *domain.User) error
}
