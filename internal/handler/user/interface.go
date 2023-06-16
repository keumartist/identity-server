package user

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	RegisterRoutes(app *fiber.App)
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}
