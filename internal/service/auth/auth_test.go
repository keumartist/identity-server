package auth_test

import (
	"art-sso/internal/domain/user"
	"art-sso/internal/service/auth"
	mocks "art-sso/internal/service/auth/mocks"
	hash "art-sso/internal/service/util"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (auth.AuthService, *mocks.MockUserRepository, *mocks.MockTokenService, *mocks.MockMailService) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenService := new(mocks.MockTokenService)
	mockMailService := new(mocks.MockMailService)
	authService := auth.NewAuthService(mockUserRepo, mockTokenService, mockMailService)
	return authService, mockUserRepo, mockTokenService, mockMailService
}

func TestAuthService(t *testing.T) {
	t.Run("Sign up with new email", func(t *testing.T) {
		authService, mockUserRepo, _, mockMailService := setup()
		email := "test@example.com"
		password := "password"

		signUpInput := auth.SignUpInput{
			Email:    email,
			Password: password,
		}

		mockUserRepo.On("GetUserByEmail", email).Return(&user.User{}, errors.New("user not found"))
		mockUserRepo.On("CreateUnverifiedUser", mock.Anything, mock.Anything).Return(nil)
		mockMailService.On("SendVerificationEmail", email, mock.Anything).Return(nil)

		_, err := authService.SignUpWithEmail(signUpInput)

		assert.Nil(t, err)
	})

	t.Run("Sign up with existing unverified email", func(t *testing.T) {
		authService, mockUserRepo, _, mockMailService := setup()
		email := "test2@example.com"
		password := "password"

		signUpInput := auth.SignUpInput{
			Email:    email,
			Password: password,
		}

		existingUser := &user.User{Email: email, EmailVerified: false}

		mockUserRepo.On("GetUserByEmail", email).Return(existingUser, nil)
		mockUserRepo.On("UpdateVerificationCode", mock.Anything, mock.Anything).Return(nil)
		mockMailService.On("SendVerificationEmail", email, mock.Anything).Return(nil)

		_, err := authService.SignUpWithEmail(signUpInput)

		assert.Nil(t, err)
	})

	t.Run("Sign in with valid credentials", func(t *testing.T) {
		authService, mockUserRepo, mockTokenService, _ := setup()
		email := "test@example.com"
		password := "password"

		signInInput := auth.SignInInput{
			Email:    email,
			Password: password,
		}

		hashedPassword, _ := hash.HashPassword(password)

		mockUserRepo.On("GetUserByEmail", email).Return(&user.User{Email: email, EmailVerified: true, Password: hashedPassword}, nil)

		mockTokenService.On("GenerateToken", mock.Anything).Return("idToken", nil)

		tokens, err := authService.SignInWithEmail(signInInput)

		assert.Nil(t, err)
		assert.NotEmpty(t, tokens.IdToken)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
	})
}
