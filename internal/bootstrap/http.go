package bootstrap

import (
	authhandler "art-sso/internal/handler/auth"
	userhandler "art-sso/internal/handler/user"
	userrepo "art-sso/internal/repository/user"
	authservice "art-sso/internal/service/auth"
	mailservice "art-sso/internal/service/mail"
	tokenservice "art-sso/internal/service/token"
	userservice "art-sso/internal/service/user"
	"log"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitHTTPServer() error {
	db, _ := gorm.Open(mysql.Open("mysql_connection_string"), &gorm.Config{})

	privateKey, secretKey, keyErr := getKeys()
	if keyErr != nil {
		return fmt.Errorf("Could not init app: %v", keyErr)
	}

	issuer, issuerErr := getIssuer()
	if issuerErr != nil {
		return fmt.Errorf("Could not init app: %v", issuerErr)
	}

	userRepo := userrepo.NewMySQLUserRepository(db)
	tokenService := tokenservice.NewTokenService(privateKey, secretKey, issuer)
	userService := userservice.NewUserService(userRepo)
	mailService := mailservice.NewMailService()
	authService := authservice.NewAuthService(userRepo, tokenService, mailService)
	userHandler := userhandler.NewUserHandler(userService)
	authHandler := authhandler.NewAuthHandler(authService)

	app := fiber.New()
	userHandler.RegisterRoutes(app)
	authHandler.RegisterRoutes(app)

	log.Fatal(app.Listen(":3000"))

	app.Listen(":3000")

	return nil
}
