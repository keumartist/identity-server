package user_test

import (
	domain "art-sso/internal/domain/user"
	service "art-sso/internal/service/user"
	"art-sso/internal/service/user/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (service.UserService, *mocks.MockUserRepository, *mocks.MockTokenService) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenService := new(mocks.MockTokenService)
	userService := service.NewUserService(mockUserRepo, mockTokenService)
	return userService, mockUserRepo, mockTokenService
}

func TestUserService(t *testing.T) {
	userService, mockUserRepo, mockTokenService := setup()

	t.Run("Create user", func(t *testing.T) {
		email := "test@example.com"
		password := "password"

		mockUserRepo.On("CreateUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.Email == email
		})).Return(nil)
		mockTokenService.On("GenerateAccessToken", "1", email).Return("test_token", nil)

		createUserInput := service.CreateUserInput{
			Email:    email,
			Password: password,
		}

		createdUser, err := userService.CreateUser(createUserInput)

		assert.Nil(t, err)

		assert.Equal(t, email, createdUser.Email)
	})

	t.Run("Get user by ID", func(t *testing.T) {
		testUser := &domain.User{ID: 1, Email: "test@example.com", Password: "password"}

		mockUserRepo.On("GetUserByID", "1").Return(testUser, nil)

		getUserByIDInput := service.GetUserByIDInput{
			ID: "1",
		}
		fetchedUser, err := userService.GetUserByID(getUserByIDInput)
		assert.Nil(t, err)
		assert.Equal(t, service.UserDomainToDto(testUser), fetchedUser)
	})
	t.Run("Get user by email", func(t *testing.T) {
		email := "test@example.com"

		testUser := &domain.User{ID: 1, Email: email, Password: "password"}

		mockUserRepo.On("GetUserByEmail", email).Return(testUser, nil)

		getUserByEmailInput := service.GetUserByEmailInput{
			Email: email,
		}

		fetchedUser, err := userService.GetUserByEmail(getUserByEmailInput)

		assert.Nil(t, err)
		assert.Equal(t, service.UserDomainToDto(testUser), fetchedUser)
	})

	t.Run("Update user", func(t *testing.T) {
		id := "1"
		newEmail := "updated@example.com"

		mockUserRepo.On("UpdateUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.ID == 1 && user.Email == newEmail
		})).Return(nil)

		updateUserInput := service.UpdateUserInput{
			ID:    id,
			Email: &newEmail,
		}

		err := userService.UpdateUser(updateUserInput)

		assert.Nil(t, err)
	})

	t.Run("Delete user", func(t *testing.T) {
		id := "1"

		mockUserRepo.On("DeleteUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.ID == 1
		})).Return(nil)

		deleteUserInput := service.DeleteUserInput{
			ID: id,
		}

		err := userService.DeleteUser(deleteUserInput)

		assert.Nil(t, err)
	})
}
