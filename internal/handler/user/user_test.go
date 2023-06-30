package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	dto "art-sso/internal/dto/user"
	mocks "art-sso/internal/handler/user/mocks"
	tokenservice "art-sso/internal/service/token"
	userservice "art-sso/internal/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	mockUserService := new(mocks.MockUserService)
	mockTokenService := new(mocks.MockTokenService)
	userHandler := NewUserHandler(mockUserService)

	app := fiber.New()
	userHandler.RegisterRoutes(app, mockTokenService)

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

	t.Run("Get my profile - failed with invalid token", func(t *testing.T) {
		userId := "test-id"
		email := "test@example.com"
		token := "valid_token"

		getUserByIdInput := userservice.GetUserByIDInput{ID: userId}
		verifyTokenInput := tokenservice.VerifyTokenInput{Token: token, TokenType: tokenservice.AccessToken}

		mockUserService.On("GetUserByID", getUserByIdInput).Return(dto.User{Email: email}, nil)
		mockTokenService.On("VerifyToken", verifyTokenInput).Return(true, userId, "", nil)

		req, _ := http.NewRequest("GET", "/users/me", nil)
		req.Header.Add("Authorization", "Bearer "+token)
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

	t.Run("Update user profile", func(t *testing.T) {
		id := "test-id"
		email := "updated@example.com"
		name := "newname"
		token := "valid_token"

		input := userservice.UpdateUserProfileInput{ID: id, Email: &email, Name: &name}
		mockUserService.On("UpdateUserProfile", input).Return(nil)
		mockTokenService.On("ValidateToken", token).Return(id, nil)

		payload := map[string]interface{}{"email": email, "name": name}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/users/"+id, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete user", func(t *testing.T) {
		id := "test-id"
		token := "valid_token"

		input := userservice.DeleteUserInput{ID: id}
		mockUserService.On("DeleteUser", input).Return(nil)
		mockTokenService.On("ValidateToken", token).Return(id, nil)

		req, _ := http.NewRequest("DELETE", "/users/"+id, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
