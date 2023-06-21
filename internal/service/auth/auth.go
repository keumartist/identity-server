package auth

import (
	"crypto/rand"
	"encoding/hex"

	"art-sso/internal/domain/user"
	errors "art-sso/internal/error"
	repository "art-sso/internal/repository/user"
	mailservice "art-sso/internal/service/mail"
	tokenservice "art-sso/internal/service/token"
	hash "art-sso/internal/service/util"
)

type AuthServiceImpl struct {
	userRepo     repository.UserRepository
	tokenService tokenservice.TokenService
	mailService  mailservice.MailService // Assuming this is another service defined somewhere else.
}

func NewAuthService(userRepo repository.UserRepository, tokenService tokenservice.TokenService, mailService mailservice.MailService) AuthService {
	return &AuthServiceImpl{
		userRepo:     userRepo,
		tokenService: tokenService,
		mailService:  mailService,
	}
}

func (s *AuthServiceImpl) SignUpWithEmail(input SignUpInput) (Tokens, error) {
	var tokens Tokens

	existingUser, err := s.userRepo.GetUserByEmail(input.Email)
	if err == nil {
		if existingUser.IsEmailVerified() {
			return tokens, errors.ErrEmailInUse
		}

		verificationCode := generateVerificationCode()
		err = s.userRepo.UpdateVerificationCode(existingUser, verificationCode)
		if err != nil {
			return tokens, errors.ErrInternal
		}

		err = s.mailService.SendVerificationEmail(input.Email, verificationCode)
		if err != nil {
			return tokens, errors.ErrSendingEmail
		}

		return tokens, nil
	}

	verificationCode := generateVerificationCode()

	hashedPassword, err := hash.HashPassword(input.Password)
	if err != nil {
		return tokens, errors.ErrInternal
	}

	newUser := &user.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	err = s.userRepo.CreateUnverifiedUser(newUser, verificationCode)
	if err != nil {
		return tokens, errors.ErrInternal
	}

	err = s.mailService.SendVerificationEmail(input.Email, verificationCode)
	if err != nil {
		return tokens, errors.ErrSendingEmail
	}

	return tokens, nil
}

func (s *AuthServiceImpl) SignInWithEmail(input SignInInput) (Tokens, error) {
	var tokens Tokens

	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		return tokens, errors.ErrInternal
	}

	if !hash.VerifyPassword(input.Password, user.Password) || !user.IsEmailVerified() {
		return tokens, errors.ErrInvalidCredentials
	}

	idToken, accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return tokens, errors.ErrInternal
	}

	tokens.IdToken = idToken
	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (s *AuthServiceImpl) generateTokens(user *user.User) (string, string, string, error) {
	idToken, err := s.tokenService.GenerateIdToken(user.IDAsString(), user.Email, 60*60*24*3)
	if err != nil {
		return "", "", "", err
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.IDAsString(), user.Email, 60*60*24*3)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.IDAsString(), user.Email, 60*60*24*7)
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
