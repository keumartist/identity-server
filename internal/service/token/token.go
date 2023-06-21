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
	AdditionalClaimsForIdToken
	jwt.StandardClaims
}

type AdditionalClaimsForIdToken struct {
	Ema string `json:"ema,omitempty"`
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

func (s *TokenServiceImpl) GenerateAccessToken(id, email string, expireAt uint) (string, error) {
	claims := Claims{
		Sub: id,
		Iss: s.issuer,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Duration(expireAt) * time.Second).Unix(),
	}

	return s.generateToken(jwt.SigningMethodHS256, claims)
}

func (s *TokenServiceImpl) GenerateRefreshToken(id, email string, expireAt uint) (string, error) {
	claims := Claims{
		Sub: id,
		Iss: s.issuer,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Duration(expireAt) * time.Second).Unix(),
	}

	return s.generateToken(jwt.SigningMethodHS256, claims)
}

func (s *TokenServiceImpl) GenerateIdToken(id, email string, expireAt uint) (string, error) {
	claims := Claims{
		Sub: id,
		Iss: s.issuer,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Duration(expireAt) * time.Second).Unix(),
		AdditionalClaimsForIdToken: AdditionalClaimsForIdToken{
			Ema: email,
		},
	}

	return s.generateToken(jwt.SigningMethodRS256, claims)

}

func (s *TokenServiceImpl) generateToken(method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	if method.Alg() == jwt.SigningMethodRS256.Alg() {
		signedToken, err := token.SignedString(s.privateKey)

		if err != nil {
			return "", errors.New("could not sign the token")
		}

		return signedToken, nil
	}

	signedToken, err := token.SignedString([]byte(s.secretKey))

	if err != nil {
		return "", errors.New("could not sign the token")
	}

	return signedToken, nil
}
