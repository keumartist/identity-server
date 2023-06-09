package user

import (
	dto "art-sso/internal/dto/user"
	customerror "art-sso/internal/error"
	mocks "art-sso/internal/handler/user/mocks"
	tokenservice "art-sso/internal/service/token"
	userservice "art-sso/internal/service/user"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	mockUserService := new(mocks.MockUserService)
	userHandler := NewUserHandler(mockUserService)

	app := fiber.New()
	mockTokenService := new(mocks.MockTokenService)
	userHandler.RegisterRoutes(app, mockTokenService)

	t.Run("Get my profile - failed with invalid token", func(t *testing.T) {
		userId := "test-id"
		email := "test@example.com"
		testToken := "test-token"

		getUserByIdInput := userservice.GetUserByIDInput{ID: userId}
		mockUserService.On("GetUserByID", getUserByIdInput).Return(dto.User{Email: email}, nil)
		mockTokenService.On("VerifyToken", tokenservice.VerifyTokenInput{Token: testToken, TokenType: tokenservice.AccessToken}).Return(true, userId, "", nil)

		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get users by email", func(t *testing.T) {
		email := "test@example.com"

		input := userservice.GetUserByEmailInput{Email: email}
		mockUserService.On("GetUserByEmail", input).Return(dto.User{Email: email}, nil)

		req, _ := http.NewRequest("GET", "/api/v1/users?email="+email, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update user profile", func(t *testing.T) {
		id := "test-id"
		name := "newname"
		testToken := "test-token"

		input := userservice.UpdateUserProfileInput{ID: id, Name: &name}
		mockUserService.On("UpdateUserProfile", input).Return(nil)
		mockTokenService.On("VerifyToken", tokenservice.VerifyTokenInput{Token: testToken, TokenType: tokenservice.AccessToken}).Return(true, id, "", nil)

		payload := map[string]interface{}{"name": name}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PUT", "/api/v1/users/me", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete user", func(t *testing.T) {
		id := "test-id"
		testToken := "test-token"

		input := userservice.DeleteUserInput{ID: id}
		mockUserService.On("DeleteUser", input).Return(nil)
		mockTokenService.On("VerifyToken", tokenservice.VerifyTokenInput{Token: testToken, TokenType: tokenservice.AccessToken}).Return(true, id, "", nil)

		req, _ := http.NewRequest("DELETE", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get my profile - failed with invalid token", func(t *testing.T) {
		testToken := "invalid-test-token"

		mockTokenService.On("VerifyToken", tokenservice.VerifyTokenInput{Token: testToken, TokenType: tokenservice.AccessToken}).Return(false, "", "", nil)

		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Get users by email - email does not exist", func(t *testing.T) {
		email := "nonexistent@example.com"

		input := userservice.GetUserByEmailInput{Email: email}
		mockUserService.On("GetUserByEmail", input).Return(dto.User{}, customerror.ErrUserNotFound)

		req, _ := http.NewRequest("GET", "/api/v1/users?email="+email, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update user profile - invalid request body", func(t *testing.T) {
		testToken := "test-token"

		invalidBody := ""
		req, _ := http.NewRequest("PUT", "/api/v1/users/me", bytes.NewBuffer([]byte(invalidBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testToken)

		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
