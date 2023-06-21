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

func (s *TokenServiceImpl) GenerateToken(input GenerateTokenInput) (string, error) {
	claims := Claims{
		Sub: input.Id,
		Iss: s.issuer,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Duration(input.ExpirationInSeconds) * time.Second).Unix(),
	}

	if input.TokenType == IdToken {
		claims.AdditionalClaimsForIdToken = AdditionalClaimsForIdToken{
			Ema: input.Email,
		}
		return s.signToken(jwt.SigningMethodRS256, claims)
	}

	return s.signToken(jwt.SigningMethodHS256, claims)
}

func (s *TokenServiceImpl) VerifyToken(input VerifyTokenInput) (bool, string, string, error) {
	var validationKey interface{}
	var claims Claims

	if input.TokenType == IdToken {
		validationKey = s.privateKey.Public()
	} else {
		validationKey = []byte(s.secretKey)
	}

	token, err := jwt.ParseWithClaims(input.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return validationKey, nil
	})

	if err != nil || !token.Valid {
		return false, "", "", errors.New("Invalid token")
	}

	return true, claims.Sub, claims.Ema, nil
}

func (s *TokenServiceImpl) signToken(method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	if method.Alg() == jwt.SigningMethodRS256.Alg() {
		return token.SignedString(s.privateKey)
	}

	return token.SignedString([]byte(s.secretKey))
}
