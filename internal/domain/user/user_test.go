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

	t.Run("create user", func(t *testing.T) {
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
}
