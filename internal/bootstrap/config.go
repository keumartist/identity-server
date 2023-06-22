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
	pemBytes, err := ioutil.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		return nil, "", fmt.Errorf("Could not read private key: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, "", fmt.Errorf("Could not parse private key: %v", err)
	}

	// Secret Key for HS256
	secretKeyBytes, err := ioutil.ReadFile(os.Getenv("SECRET_KEY_PATH"))
	if err != nil {
		return nil, "", fmt.Errorf("Could not read secret key: %v", err)
	}
	secretKey := string(secretKeyBytes)

	return privateKey, secretKey, nil
}

func getIssuer() (string, error) {
	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		return "", errors.New("Issuer not found")
	}

	return issuer, nil
}
