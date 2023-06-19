package user

import (
	domain "art-sso/internal/domain/user"
	dto "art-sso/internal/dto/user"
	repository "art-sso/internal/repository/user"
	util "art-sso/internal/service/util"
	"fmt"
	"strconv"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) CreateUser(input CreateUserInput) (dto.User, error) {
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return dto.User{}, err
	}

	user := &domain.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return dto.User{}, err
	}

	return UserDomainToDto(user), nil
}

func (s *UserServiceImpl) GetUserByID(input GetUserByIDInput) (dto.User, error) {
	user, err := s.repo.GetUserByID(input.ID)

	if err != nil {
		return dto.User{}, err
	}

	return UserDomainToDto(user), nil
}

func (s *UserServiceImpl) GetUserByEmail(input GetUserByEmailInput) (dto.User, error) {
	user, err := s.repo.GetUserByEmail(input.Email)

	if err != nil {
		return dto.User{}, err
	}

	return UserDomainToDto(user), nil
}

func (s *UserServiceImpl) UpdateUserProfile(input UpdateUserProfileInput) error {
	id, err := strconv.ParseUint(input.ID, 10, 32)

	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}

	user := &domain.User{
		ID:   uint(id),
		Name: *input.Name,
	}

	err = s.repo.UpdateUserProfile(user)
	if err != nil {
		return fmt.Errorf("failed to update user profile in repository: %v", err)
	}

	return nil
}

func (s *UserServiceImpl) DeleteUser(input DeleteUserInput) error {
	id, err := strconv.ParseUint(input.ID, 10, 32)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID: uint(id),
	}

	return s.repo.DeleteUser(user)
}
