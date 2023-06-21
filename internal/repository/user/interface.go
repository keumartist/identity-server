package user

import domain "art-sso/internal/domain/user"

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	UpdateUserProfile(user *domain.User) error
	UpdateVerificationCode(user *domain.User, verficationCode string) error
	DeleteUser(user *domain.User) error
	CreateUnverifiedUser(user *domain.User, verificationCode string) error
	VerifyUser(email, verificationCode string) error
}
