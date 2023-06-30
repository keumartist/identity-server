package middleware

import (
	customerror "art-sso/internal/error"
	"art-sso/internal/service/token"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TokenValidationMiddleware(s token.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenHeader := c.Get("Authorization")
		if tokenHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(customerror.ErrUnauthorized)
		}

		splitToken := strings.Split(tokenHeader, "Bearer ")
		if len(splitToken) != 2 {
			return c.Status(http.StatusUnauthorized).JSON(customerror.ErrUnauthorized)
		}
		authToken := splitToken[1]

		isValid, userId, _, err := s.VerifyToken(token.VerifyTokenInput{
			Token:     authToken,
			TokenType: token.AccessToken,
		})
		if err != nil || !isValid {
			return c.Status(http.StatusUnauthorized).JSON(customerror.ErrUnauthorized)
		}

		c.Locals("userId", userId)

		return c.Next()
	}
}
