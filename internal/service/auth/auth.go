package auth

import (
	domain "art-sso/internal/domain/user"
	customerror "art-sso/internal/error"
	repository "art-sso/internal/repository/user"
	mailservice "art-sso/internal/service/mail"
	tokenservice "art-sso/internal/service/token"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	userRepo     repository.UserRepository
	tokenService tokenservice.TokenService
	mailService  mailservice.MailService
	oauthConfig  *oauth2.Config
}

func NewAuthService(userRepo repository.UserRepository, tokenService tokenservice.TokenService, mailService mailservice.MailService) AuthService {
	oauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	return &AuthServiceImpl{
		userRepo:     userRepo,
		tokenService: tokenService,
		mailService:  mailService,
		oauthConfig:  oauthConfig,
	}
}

func (s *AuthServiceImpl) generateTokens(user *domain.User) (string, string, string, error) {
	idToken, err := s.tokenService.GenerateToken(tokenservice.GenerateTokenInput{Id: user.IDAsString(), Email: user.Email, ExpirationInSeconds: 60 * 60 * 24 * 3, TokenType: tokenservice.IdToken})
	if err != nil {
		return "", "", "", err
	}

	accessToken, err := s.tokenService.GenerateToken(tokenservice.GenerateTokenInput{Id: user.IDAsString(), Email: user.Email, ExpirationInSeconds: 60 * 60 * 24 * 3, TokenType: tokenservice.AccessToken})
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := s.tokenService.GenerateToken(tokenservice.GenerateTokenInput{Id: user.IDAsString(), Email: user.Email, ExpirationInSeconds: 60 * 60 * 24 * 7, TokenType: tokenservice.RefreshToken})
	if err != nil {
		return "", "", "", err
	}

	return idToken, accessToken, refreshToken, nil
}

func generateVerificationCodeWithExpireTime(timeInMills uint) (string, time.Time) {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	expireAt := time.Now().Add(time.Duration(timeInMills) * time.Second)

	return hex.EncodeToString(bytes), expireAt
}

func (s *AuthServiceImpl) RefreshAccessToken(input RefreshAccessTokenInput) (Tokens, error) {
	id, err := s.tokenService.GetUserIDFromToken(tokenservice.GetUserIDFromTokenInput{Token: input.Token, TokenType: tokenservice.RefreshToken})

	if err != nil {
		return Tokens{}, customerror.ErrUnauthorized
	}

	user, err := s.userRepo.GetUserByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Tokens{}, customerror.ErrUnauthorized
		}

		return Tokens{}, customerror.ErrInternal
	}

	if user.RefreshToken != input.Token {
		return Tokens{}, customerror.ErrUnauthorized
	}

	_, accessToken, _, err := s.generateTokens(user)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{AccessToken: accessToken}, nil
}
