package service_test

import (
	"testing"

	domain "art-sso/internal/domain/user"
	repository "art-sso/internal/repository/user"
	service "art-sso/internal/service/user"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&domain.User{})

	userRepository := repository.NewMySQLUserRepository(db)
	userService := service.NewUserService(userRepository)

	t.Run("Create user", func(t *testing.T) {
		email := "user@example.com"
		password := "password"

		user, err := userService.CreateUser(email, password)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
	})

	t.Run("Get user by ID", func(t *testing.T) {
		id := "1"

		user, err := userService.GetUserByID(id)
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", user.Email)
	})

	t.Run("Get user by email", func(t *testing.T) {
		email := "user@example.com"

		user, err := userService.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
	})

	t.Run("Update user", func(t *testing.T) {
		user, err := userService.GetUserByID("1")
		assert.NoError(t, err)

		user.Email = "updated@example.com"
		err = userService.UpdateUser(user)
		assert.NoError(t, err)

		updatedUser, err := userService.GetUserByID("1")
		assert.NoError(t, err)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	t.Run("Delete user", func(t *testing.T) {
		user, err := userService.GetUserByID("1")
		assert.NoError(t, err)
		err = userService.DeleteUser(user)
		assert.NoError(t, err)

		user, err = userService.GetUserByID("1")
		assert.Nil(t, user)
	})
}
