package auth_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	authhandler "art-sso/internal/handler/auth"
	mocks "art-sso/internal/handler/auth/mocks"
	authservice "art-sso/internal/service/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setUp() (*mocks.MockAuthService, *fiber.App) {
	mockAuthService := new(mocks.MockAuthService)
	handler := authhandler.NewAuthHandler(mockAuthService)

	app := fiber.New()
	handler.RegisterRoutes(app)

	return mockAuthService, app
}

func TestSignUpWithEmail(t *testing.T) {
	t.Run("Successful sign up", func(t *testing.T) {
		mockAuthService, app := setUp()

		mockAuthService.On("SignUpWithEmail", authservice.SignUpInput{
			Email:    "test@example.com",
			Password: "password",
		}).Return(authservice.Tokens{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}, nil)

		req, _ := http.NewRequest("POST", "/signup", strings.NewReader(`{"email":"test@example.com", "password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Failed sign up due to email in use", func(t *testing.T) {
		mockAuthService, app := setUp()

		mockAuthService.On("SignUpWithEmail", authservice.SignUpInput{
			Email:    "test@example.com",
			Password: "password",
		}).Return(authservice.Tokens{}, errors.New("email already in use"))

		req, _ := http.NewRequest("POST", "/signup", strings.NewReader(`{"email":"test@example.com", "password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
