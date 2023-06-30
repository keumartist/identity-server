package auth

import (
	"crypto/rand"
	"encoding/hex"

	"art-sso/internal/domain/user"
	customerror "art-sso/internal/error"
	tokenservice "art-sso/internal/service/token"
	hash "art-sso/internal/service/util"
	"errors"

	"gorm.io/gorm"
)

func (s *AuthServiceImpl) SignUpWithEmail(input SignUpInput) (string, error) {
	existingUser, err := s.userRepo.GetUserByEmail(input.Email)
	if err == nil {
		return s.updateExistingUser(existingUser)
	}

	return s.createNewUser(input.Email, input.Password)
}

func (s *AuthServiceImpl) SignInWithEmail(input SignInInput) (Tokens, error) {
	var tokens Tokens

	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tokens, customerror.ErrInvalidCredentials
		}

		return tokens, customerror.ErrInternal
	}

	if !hash.VerifyPassword(input.Password, user.Password) || !user.IsEmailVerified() {
		return tokens, customerror.ErrInvalidCredentials
	}

	idToken, accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return tokens, customerror.ErrInternal
	}

	tokens.IdToken = idToken
	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (s *AuthServiceImpl) updateExistingUser(existingUser *user.User) (string, error) {
	if existingUser.IsEmailVerified() {
		return "", customerror.ErrEmailInUse
	}

	verificationCode := generateVerificationCode()
	err := s.userRepo.UpdateVerificationCode(existingUser, verificationCode)
	if err != nil {
		return "", customerror.ErrInternal
	}

	err = s.mailService.SendVerificationEmail(existingUser.Email, verificationCode)
	if err != nil {
		return "", customerror.ErrSendingEmail
	}

	return "Verification code was sent to user email", nil
}

func (s *AuthServiceImpl) createNewUser(email, password string) (string, error) {
	verificationCode := generateVerificationCode()

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return "", customerror.ErrInternal
	}

	newUser := &user.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = s.userRepo.CreateUnverifiedUser(newUser, verificationCode)
	if err != nil {
		return "", customerror.ErrInternal
	}

	err = s.mailService.SendVerificationEmail(email, verificationCode)
	if err != nil {
		return "", customerror.ErrSendingEmail
	}

	return "Verification code was sent to user email", nil
}

func (s *AuthServiceImpl) generateTokens(user *user.User) (string, string, string, error) {
	idToken, err := s.tokenService.GenerateToken(tokenservice.GenerateTokenInput{Id: user.IDAsString(), Email: user.Email, ExpirationInSeconds: 60 * 60 * 24 * 3, TokenType: tokenservice.IdToken})
	if err != nil {
		return "", "", "", err
	}

	accessToken, err := s.tokenService.GenerateToken(tokenservice.GenerateTokenInput{Id: user.IDAsString(), Email: user.Email, ExpirationInSeconds: 60 * 60 * 24 * 3, TokenType: tokenservice.AccessToken})
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := s.tokenService.GenerateToken(tokenservice.GenerateTokenInput{user.IDAsString(), user.Email, 60 * 60 * 24 * 7, tokenservice.RefreshToken})
	if err != nil {
		return "", "", "", err
	}

	return idToken, accessToken, refreshToken, nil
}

func generateVerificationCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
