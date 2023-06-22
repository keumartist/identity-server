package bootstrap

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
)

func getKeys() (*rsa.PrivateKey, string, error) {
	// Private Key for RS256
	pemBytes, err := ioutil.ReadFile(os.Getenv("PRIVATE_KEY_PATH_FOR_RS256_ID_TOKEN"))
	if err != nil {
		return nil, "", fmt.Errorf("Could not read private key: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, "", fmt.Errorf("Could not parse private key: %v", err)
	}

	// Secret Key for HS256
	secretKey := os.Getenv("SECRET_KEY_FOR_HS256_ACCESS_TOKEN")
	if secretKey == "" {
		return nil, "", errors.New("Secret key not found")
	}

	return privateKey, secretKey, nil
}

func getIssuer() (string, error) {
	issuer := os.Getenv("ISSUER_FOR_TOKEN")
	if issuer == "" {
		return "", errors.New("Issuer not found")
	}

	return issuer, nil
}
