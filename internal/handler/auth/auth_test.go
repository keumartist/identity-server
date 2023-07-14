package auth_test

import (
	"net/http"
	"strings"
	"testing"

	customerror "art-sso/internal/error"
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
		}).Return("Verification code was sent to email", nil)

		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(`{"email":"test@example.com", "password":"password"}`))
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
		}).Return("", customerror.ErrEmailInUse)

		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(`{"email":"test@example.com", "password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestVerifyEmailCode(t *testing.T) {
	t.Run("Failed email verification due to invalid code", func(t *testing.T) {
		mockAuthService, app := setUp()

		mockAuthService.On("VerifyEmailCode", authservice.VerifyEmailCodeInput{
			Email: "test@example.com",
			Code:  "code",
		}).Return(customerror.ErrInvalidVerificationCode)

		req, _ := http.NewRequest("POST", "/api/v1/auth/verification", strings.NewReader(`{"email":"test@example.com", "code":"code"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Successful email verification", func(t *testing.T) {
		mockAuthService, app := setUp()

		mockAuthService.On("VerifyEmailCode", authservice.VerifyEmailCodeInput{
			Email: "test@example.com",
			Code:  "code",
		}).Return(nil)

		req, _ := http.NewRequest("POST", "/api/v1/auth/verification", strings.NewReader(`{"email":"test@example.com", "code":"code"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
