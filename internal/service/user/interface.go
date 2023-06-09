package user

import domain "art-sso/internal/domain/user"

type UserService interface {
	CreateUser(email, password string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(user *domain.User) error
}
