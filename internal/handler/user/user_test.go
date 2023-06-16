package user

import (
	domain "art-sso/internal/domain/user"
	repository "art-sso/internal/repository/user"
	service "art-sso/internal/service/user"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserHandler(t *testing.T) {
	// Initialize the database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&domain.User{})

	userRepository := repository.NewMySQLUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := NewUserHandler(userService)

	app := fiber.New()
	userHandler.RegisterRoutes(app)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
	})

	t.Run("Create user", func(t *testing.T) {
		payload := map[string]string{
			"email":    "test@example.com",
			"password": "password",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		req, _ = http.NewRequest("GET", "/users?email=test@example.com", nil)
		resp, err = app.Test(req, -1)

		var user domain.User
		json.NewDecoder(resp.Body).Decode(&user)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("Update user", func(t *testing.T) {
		payload := map[string]interface{}{
			"ID":    1,
			"Email": "updated@example.com",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Fetch again to check
		req, _ = http.NewRequest("GET", "/users/1", nil)
		resp, err = app.Test(req, -1)
		var updatedUser domain.User
		json.NewDecoder(resp.Body).Decode(&updatedUser)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	t.Run("Delete user", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Fetch again to check if deleted
		req, _ = http.NewRequest("GET", "/users/1", nil)
		resp, err = app.Test(req, -1)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
