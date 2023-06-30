package auth

import (
	repository "art-sso/internal/repository/user"
	mailservice "art-sso/internal/service/mail"
	tokenservice "art-sso/internal/service/token"
	"os"

	"golang.org/x/oauth2"
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
