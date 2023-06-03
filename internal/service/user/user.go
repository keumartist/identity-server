package service

import (
	domain "art-sso/internal/domain/user"
	repository "art-sso/internal/repository/user"
	util "art-sso/internal/service/util"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(email, password string) (*domain.User, error) {
	salt := util.GenerateSaltForPassword(password)

	user := &domain.User{
		Email:    email,
		Password: password,
		Salt:     salt,
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(user *domain.User) error {
	return s.repo.DeleteUser(user)
}
