package user

import (
	domain "art-sso/internal/domain/user"
	repository "art-sso/internal/repository/user"
	tokenservice "art-sso/internal/service/token"
	util "art-sso/internal/service/util"
)

type UserServiceImpl struct {
	repo     repository.UserRepository
	tokenSvc tokenservice.TokenService
}

func NewUserService(repo repository.UserRepository, tokenSvc tokenservice.TokenService) *UserServiceImpl {
	return &UserServiceImpl{
		repo:     repo,
		tokenSvc: tokenSvc,
	}
}

func (s *UserServiceImpl) CreateUser(email, password string) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = s.repo.CreateUser(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetUserByID(id string) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserServiceImpl) GetUserByEmail(email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserServiceImpl) UpdateUser(user *domain.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserServiceImpl) DeleteUser(user *domain.User) error {
	return s.repo.DeleteUser(user)
}
