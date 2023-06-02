package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	domain "art-sso/internal/domain/user"
)

func TestUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.User{})
	assert.NoError(t, err)

	t.Run("Create user", func(t *testing.T) {
		user := domain.User{
			Email:         "user@example.com",
			Salt:          "randomsalt",
			EmailVerified: false,
		}

		result := db.Create(&user)
		assert.NoError(t, result.Error)
		assert.NotNil(t, user.ID)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
		assert.Equal(t, "user@example.com", user.Email)
		assert.Equal(t, "randomsalt", user.Salt)
		assert.False(t, user.EmailVerified)
	})

	t.Run("Change password", func(t *testing.T) {
		user := domain.User{
			Email:         "user2@example.com",
			Password:      "oldpassword",
			Salt:          "randomsalt",
			EmailVerified: false,
		}

		result := db.Create(&user)
		assert.NoError(t, result.Error)

		err := user.ChangePassword("oldpassword", "newpassword")
		assert.NoError(t, err)
		assert.Equal(t, "newpassword", user.Password)
	})

	t.Run("Verify email", func(t *testing.T) {
		code := "verificationcode"
		user := domain.User{
			Email:            "user3@example.com",
			VerificationCode: &code,
			EmailVerified:    false,
		}

		result := db.Create(&user)
		assert.NoError(t, result.Error)

		err := user.VerifyEmail("verificationcode")
		assert.NoError(t, err)
		assert.True(t, user.EmailVerified)
	})
}
