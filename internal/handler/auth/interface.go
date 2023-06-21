package auth

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	RegisterRoutes(app *fiber.App)
	SignUpWithEmail(c *fiber.Ctx) error
	SignInWithEmail(c *fiber.Ctx) error
}
