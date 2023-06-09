package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	domain "art-sso/internal/domain/user"
)

func TestMySQLUserRepository(t *testing.T) {
	// MySQL Repository 테스트이지만, Unit Test 에서는 일단은 SQLite를 사용한다. 현재 메소드는 모두 서로 호환되는 상황.
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.SocialConnection{})
	if err != nil {
		t.Fatal(err)
	}

	userRepo := NewMySQLUserRepository(db)

	t.Run("Create user", func(t *testing.T) {
		user := domain.User{
			Email:         "user@example.com",
			EmailVerified: false,
		}

		err := userRepo.CreateUser(&user)

		assert.NoError(t, err)
	})

	t.Run("Get user by ID", func(t *testing.T) {
		user, err := userRepo.GetUserByID("1")
		t.Log(user)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("Get user by email", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail("user@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("Update user", func(t *testing.T) {
		user, err := userRepo.GetUserByID("1")

		assert.NoError(t, err)

		user.EmailVerified = true
		err = userRepo.UpdateUser(user)

		assert.NoError(t, err)
	})

	t.Run("Update user profile", func(t *testing.T) {
		user, err := userRepo.GetUserByID("1")
		assert.NoError(t, err)

		user.Name = "newusername"

		err = userRepo.UpdateUserProfile(user)
		assert.NoError(t, err)

		updatedUser, err := userRepo.GetUserByID("1")
		assert.NoError(t, err)
		assert.Equal(t, "newusername", updatedUser.Name)
	})

	t.Run("Delete user", func(t *testing.T) {
		user, err := userRepo.GetUserByID("1")

		assert.NoError(t, err)

		err = userRepo.DeleteUser(user)

		assert.NoError(t, err)
	})

	t.Run("Create unverified user", func(t *testing.T) {
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")

		verificationCode := "123456"
		expireAt := time.Now().Add(time.Duration(180) * time.Second)
		user := domain.User{
			Email: "unverified@example.com",
		}

		err := userRepo.CreateUnverifiedUser(&user, verificationCode, expireAt)

		assert.NoError(t, err)

		createdUser, err := userRepo.GetUserByEmail("unverified@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.False(t, createdUser.EmailVerified)
		assert.Equal(t, verificationCode, *createdUser.VerificationCode)
	})

	t.Run("Verify user email", func(t *testing.T) {
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")

		email := "unverified@example.com"
		user := domain.User{
			Email: email,
		}

		err := userRepo.VerifyUserEmail(&user)

		assert.NoError(t, err)

		verifiedUser, err := userRepo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.NotNil(t, verifiedUser)
		assert.True(t, verifiedUser.EmailVerified)

		assert.Equal(t, "", *verifiedUser.VerificationCode)
	})
}
