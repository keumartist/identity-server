package user

import (
	tokenservice "art-sso/internal/service/token"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	RegisterRoutes(app *fiber.App, tokenService tokenservice.TokenService)
	GetMe(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	UpdateMeUserProfile(c *fiber.Ctx) error
	DeleteMeUser(c *fiber.Ctx) error
}
