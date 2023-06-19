package user_test

import (
	domain "art-sso/internal/domain/user"
	user "art-sso/internal/service/user"
	mocks "art-sso/internal/service/user/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (user.UserService, *mocks.MockUserRepository) {
	mockUserRepo := new(mocks.MockUserRepository)
	userService := user.NewUserService(mockUserRepo)
	return userService, mockUserRepo
}

func TestUserService(t *testing.T) {
	userService, mockUserRepo := setup()

	t.Run("Create user", func(t *testing.T) {
		email := "test@example.com"
		password := "password"

		mockUserRepo.On("CreateUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.Email == email
		})).Return(nil)

		createUserInput := user.CreateUserInput{
			Email:    email,
			Password: password,
		}

		createdUser, err := userService.CreateUser(createUserInput)

		assert.Nil(t, err)
		assert.Equal(t, email, createdUser.Email)
	})

	t.Run("Get user by ID", func(t *testing.T) {
		id := "1"
		email := "test@example.com"

		getUserByIdInput := user.GetUserByIDInput{
			ID: id,
		}

		mockUserRepo.On("GetUserByID", id).Return(&domain.User{Email: email}, nil)

		fetchedUser, err := userService.GetUserByID(getUserByIdInput)

		assert.Nil(t, err)
		assert.Equal(t, email, fetchedUser.Email)
	})

	t.Run("Get user by email", func(t *testing.T) {
		email := "test@example.com"

		getUserByEmailInput := user.GetUserByEmailInput{
			Email: email,
		}

		mockUserRepo.On("GetUserByEmail", email).Return(&domain.User{Email: email}, nil)

		fetchedUser, err := userService.GetUserByEmail(getUserByEmailInput)

		assert.Nil(t, err)
		assert.Equal(t, email, fetchedUser.Email)
	})

	t.Run("Update user profile", func(t *testing.T) {
		newName := "new name"

		updateUserInput := user.UpdateUserProfileInput{
			ID:   "3",
			Name: &newName,
		}

		mockUserRepo.On("UpdateUserProfile", mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == newName
		})).Return(nil)

		err := userService.UpdateUserProfile(updateUserInput)

		assert.NoError(t, err)
	})

	t.Run("Delete user", func(t *testing.T) {
		userToDelete := &domain.User{
			ID:    3,
			Email: "test@example.com",
		}

		deleteUserInput := user.DeleteUserInput{
			ID: "3",
		}

		mockUserRepo.On("DeleteUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.ID == userToDelete.ID
		})).Return(nil)

		err := userService.DeleteUser(deleteUserInput)

		assert.Nil(t, err)
	})
}
