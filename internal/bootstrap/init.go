package bootstrap

import (
	handler "art-sso/internal/handler/user"
	repo "art-sso/internal/repository/user"
	tokenservice "art-sso/internal/service/token"
	userservice "art-sso/internal/service/user"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitApp() *fiber.App {
	db, _ := gorm.Open(mysql.Open("mysql_connection_string"), &gorm.Config{})

	privateKey, secretKey, keyErr := getKeys()
	if keyErr != nil {
		fmt.Printf("Could not init app: %v", keyErr)
		return nil
	}

	issuer, issuerErr := getIssuer()
	if issuerErr != nil {
		fmt.Printf("Could not init app: %v", issuerErr)
		return nil
	}

	userRepo := repo.NewMySQLUserRepository(db)
	tokenService := tokenservice.NewTokenService(privateKey, secretKey, issuer)
	userService := userservice.NewUserService(userRepo, tokenService)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()
	userHandler.RegisterRoutes(app)

	return app
}

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
