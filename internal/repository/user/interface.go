package user

import (
	domain "art-sso/internal/domain/user"
	"time"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	UpdateUserProfile(user *domain.User) error
	UpdateVerificationCode(user *domain.User, verificationCode string, expireAt time.Time) error
	UpdateRefreshToken(user *domain.User, refreshToken string) error
	DeleteUser(user *domain.User) error
	CreateUnverifiedUser(user *domain.User, verificationCode string, expireAt time.Time) error
	VerifyUserEmail(user *domain.User) error
}
