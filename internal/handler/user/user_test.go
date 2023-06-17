package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	dto "art-sso/internal/dto/user"
	mock "art-sso/internal/handler/user/mocks"
	userservice "art-sso/internal/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	mockUserService := new(mock.MockUserService)
	userHandler := NewUserHandler(mockUserService)

	app := fiber.New()
	userHandler.RegisterRoutes(app)

	t.Run("Create user", func(t *testing.T) {
		email := "test@example.com"
		password := "password"

		input := userservice.CreateUserInput{Email: email, Password: password}
		mockUserService.On("CreateUser", input).Return(dto.User{Email: email}, nil)

		payload := map[string]string{
			"email":    email,
			"password": password,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Get user by ID", func(t *testing.T) {
		id := "test-id"
		email := "test@example.com"

		input := userservice.GetUserByIDInput{ID: id}
		mockUserService.On("GetUserByID", input).Return(dto.User{Email: email}, nil)

		req, _ := http.NewRequest("GET", "/users/"+id, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get users by email", func(t *testing.T) {
		email := "test@example.com"

		input := userservice.GetUserByEmailInput{Email: email}
		mockUserService.On("GetUserByEmail", input).Return(dto.User{Email: email}, nil)

		req, _ := http.NewRequest("GET", "/users?email="+email, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update user", func(t *testing.T) {
		id := "test-id"
		email := "updated@example.com"
		password := ""

		input := userservice.UpdateUserInput{ID: id, Email: &email, Password: &password}
		mockUserService.On("UpdateUser", input).Return(nil)

		payload := map[string]interface{}{"id": id, "email": email}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/users/"+id, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete user", func(t *testing.T) {
		id := "test-id"

		input := userservice.DeleteUserInput{ID: id}
		mockUserService.On("DeleteUser", input).Return(nil)

		req, _ := http.NewRequest("DELETE", "/users/"+id, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
