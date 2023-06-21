package token_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"art-sso/internal/service/token"
	tokenservice "art-sso/internal/service/token" // your project path to token package

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func setup() (tokenservice.TokenService, *rsa.PrivateKey, string) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	secretKey := "mySecretKey"
	issuer := "testIssuer"

	tokenService := tokenservice.NewTokenService(privateKey, secretKey, issuer)
	return tokenService, privateKey, secretKey
}

func TestTokenService(t *testing.T) {
	tokenService, privateKey, secretKey := setup()

	t.Run("Generate access token", func(t *testing.T) {
		expireAt := uint(3600)
		tokenString, err := tokenService.GenerateAccessToken("1", "test1@example.com", expireAt)
		assert.Nil(t, err)

		token, err := jwt.ParseWithClaims(tokenString, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		assert.Nil(t, err)

		claims, ok := token.Claims.(*tokenservice.Claims)
		assert.True(t, ok)

		expiration := time.Unix(claims.Exp, 0)
		expectedExpiration := time.Now().Add(time.Duration(expireAt) * time.Second)
		diffInSeconds := int64(expectedExpiration.Sub(expiration).Seconds())
		fmt.Println(diffInSeconds)
		assert.True(t, diffInSeconds >= 0 && diffInSeconds <= 1)
	})

	t.Run("Generate refresh token", func(t *testing.T) {
		expireAt := uint(7200)
		tokenString, err := tokenService.GenerateRefreshToken("2", "test2@example.com", expireAt)
		assert.Nil(t, err)

		token, err := jwt.ParseWithClaims(tokenString, &tokenservice.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		assert.Nil(t, err)
		claims, ok := token.Claims.(*tokenservice.Claims)
		assert.True(t, ok)
		assert.Equal(t, "2", claims.Sub)
		assert.Equal(t, "testIssuer", claims.Iss)

		expiration := time.Unix(claims.Exp, 0)
		expectedExpiration := time.Now().Add(time.Duration(expireAt) * time.Second)
		diffInSeconds := int64(expectedExpiration.Sub(expiration).Seconds())
		fmt.Println(diffInSeconds)
		assert.True(t, diffInSeconds >= 0 && diffInSeconds <= 1)
	})

	t.Run("Generate ID token", func(t *testing.T) {
		expireAt := uint(3600)
		tokenString, err := tokenService.GenerateIdToken("3", "test3@example.com", expireAt)
		assert.Nil(t, err)

		token, err := jwt.ParseWithClaims(tokenString, &tokenservice.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return &privateKey.PublicKey, nil
		})

		assert.Nil(t, err)
		claims, ok := token.Claims.(*tokenservice.Claims)
		assert.True(t, ok)
		assert.Equal(t, "3", claims.Sub)
		assert.Equal(t, "testIssuer", claims.Iss)
		assert.Equal(t, "test3@example.com", claims.Ema)

		expiration := time.Unix(claims.Exp, 0)
		expectedExpiration := time.Now().Add(time.Duration(expireAt) * time.Second)
		diffInSeconds := int64(expectedExpiration.Sub(expiration).Seconds())
		fmt.Println(diffInSeconds)
		assert.True(t, diffInSeconds >= 0 && diffInSeconds <= 1)
	})
}
