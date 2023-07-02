package auth

import (
	customerror "art-sso/internal/error"
	"context"
	"errors"

	domain "art-sso/internal/domain/user"

	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"gorm.io/gorm"
)

func (s *AuthServiceImpl) SignInWithGoogle(input SignInWithGoogleInput) (Tokens, error) {
	var tokens Tokens

	token, err := s.oauthConfig.Exchange(context.Background(), input.Code)
	if err != nil {
		return tokens, customerror.ErrInvalidCredentials
	}

	peopleService, err := people.NewService(context.Background(), option.WithTokenSource(s.oauthConfig.TokenSource(context.Background(), token)))
	if err != nil {
		return tokens, customerror.ErrInternal
	}

	person, err := peopleService.People.Get("people/me").PersonFields("names,emailAddresses").Do()
	if err != nil {
		return tokens, customerror.ErrInternal
	}

	user, err := s.userRepo.GetUserByEmail(person.EmailAddresses[0].Value)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return tokens, customerror.ErrInternal
	}

	if user == nil {
		user = &domain.User{
			Email: person.EmailAddresses[0].Value,
			Name:  person.Names[0].DisplayName,
		}

		socialConnection := &domain.SocialConnection{
			SocialProviderID: domain.GoogleProviderID,
			SocialMediaID:    person.ResourceName,
		}

		if err := socialConnection.Connect(user, &domain.SocialProvider{ID: domain.GoogleProviderID, Name: "Google"}); err != nil {
			return tokens, customerror.ErrInternal
		}

		if err := s.userRepo.CreateUser(user); err != nil {
			return tokens, customerror.ErrInternal
		}
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
