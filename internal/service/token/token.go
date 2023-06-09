package token

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Sub string `json:"sub"`
	Aud string `json:"aud,omitempty"`
	Iss string `json:"iss,omitempty"`
	Iat int64  `json:"iat,omitempty"`
	Exp int64  `json:"exp,omitempty"`
	jwt.StandardClaims
}

type TokenServiceImpl struct {
	privateKey *rsa.PrivateKey
	secretKey  string
	issuer     string
}

func NewTokenService(privateKey *rsa.PrivateKey, secretKey, issuer string) TokenService {
	return &TokenServiceImpl{
		privateKey: privateKey,
		secretKey:  secretKey,
		issuer:     issuer,
	}
}

func (s *TokenServiceImpl) GenerateAccessToken(id, email string) (string, error) {
	claims := Claims{
		Sub: id,
		Iss: s.issuer,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	return s.generateToken(jwt.SigningMethodHS256, claims)
}

func (s *TokenServiceImpl) GenerateRefreshToken(id, email string) (string, error) {
	claims := Claims{
		Sub: id,
		Iss: s.issuer,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24 * 14).Unix(),
	}

	return s.generateToken(jwt.SigningMethodHS256, claims)
}

func (s *TokenServiceImpl) generateToken(method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", errors.New("could not sign the token")
	}

	return signedToken, nil
}
