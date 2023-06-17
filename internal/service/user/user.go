package user

import (
	domain "art-sso/internal/domain/user"
	dto "art-sso/internal/dto/user"
	userrepository "art-sso/internal/repository/user"
	tokenservice "art-sso/internal/service/token"
	util "art-sso/internal/service/util"
	"strconv"
)

type UserServiceImpl struct {
	repo     userrepository.UserRepository
	tokenSvc tokenservice.TokenService
}

func NewUserService(repo userrepository.UserRepository, tokenSvc tokenservice.TokenService) UserService {
	return &UserServiceImpl{
		repo:     repo,
		tokenSvc: tokenSvc,
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

func (s *UserServiceImpl) UpdateUser(input UpdateUserInput) error {
	id, err := strconv.ParseUint(input.ID, 10, 32)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID: uint(id),
	}

	// Check if Email is provided
	if input.Email != nil {
		user.Email = *input.Email
	}

	// Check if Password is provided
	if input.Password != nil {
		hashedPassword, err := util.HashPassword(*input.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	err = s.repo.UpdateUser(user)

	return err
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
