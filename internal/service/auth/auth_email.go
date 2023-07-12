package auth

import (
	"time"

	"art-sso/internal/domain/user"
	customerror "art-sso/internal/error"
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

	err = s.userRepo.UpdateRefreshToken(user, refreshToken)
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

	verificationCode, expireAt := generateVerificationCodeWithExpireTime(180)

	err := s.userRepo.UpdateVerificationCode(existingUser, verificationCode, expireAt)
	if err != nil {
		return "", customerror.ErrInternal
	}

	err = s.mailService.SendVerificationEmail(existingUser.Email, verificationCode)
	if err != nil {
		return "", customerror.ErrSendingEmail
	}

	return "Verification code was sent to user email", nil
}

func (s *AuthServiceImpl) VerifyEmailCode(input VerifyEmailCodeInput) error {
	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		return customerror.ErrUserNotFound
	}

	if *user.VerificationCode != input.Code {
		return customerror.ErrInvalidVerificationCode
	}

	if time.Now().After(*user.VerificationCodeExpireAt) {
		return customerror.ErrInvalidVerificationCode
	}

	err = s.userRepo.VerifyUserEmail(user)
	if err != nil {
		return customerror.ErrInvalidVerificationCode
	}

	return nil
}

func (s *AuthServiceImpl) createNewUser(email, password string) (string, error) {
	verificationCode, expireAt := generateVerificationCodeWithExpireTime(180)

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return "", customerror.ErrInternal
	}

	newUser := &user.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = s.userRepo.CreateUnverifiedUser(newUser, verificationCode, expireAt)
	if err != nil {
		return "", customerror.ErrInternal
	}

	err = s.mailService.SendVerificationEmail(email, verificationCode)
	if err != nil {
		return "", customerror.ErrSendingEmail
	}

	return "Verification code was sent to user email", nil
}
