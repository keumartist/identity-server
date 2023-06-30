package user

import (
	"net/http"

	tokenservice "art-sso/internal/service/token"
	userservice "art-sso/internal/service/user"

	middleware "art-sso/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerImpl struct {
	service userservice.UserService
}

func NewUserHandler(service userservice.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		service: service,
	}
}

func (h *UserHandlerImpl) RegisterRoutes(app *fiber.App, tokenService tokenservice.TokenService) {
	app.Post("/users", h.CreateUser)
	app.Get("/users/me", middleware.TokenValidationMiddleware(tokenService), h.GetMe)
	app.Get("/users", h.GetUsers)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)
}

func (h *UserHandlerImpl) CreateUser(c *fiber.Ctx) error {
	var requestBody CreateUserRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	input := userservice.CreateUserInput{Email: requestBody.Email, Password: requestBody.Password}
	createdUser, err := h.service.CreateUser(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("User not found")
	}

	return c.Status(http.StatusCreated).JSON(createdUser)
}

func (h *UserHandlerImpl) GetMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).SendString("Not authorized")
	}

	input := userservice.GetUserByIDInput{ID: userId}

	user, err := h.service.GetUserByID(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(user)
}

func (h *UserHandlerImpl) GetUsers(c *fiber.Ctx) error {
	email := c.Query("email")
	input := userservice.GetUserByEmailInput{Email: email}

	user, err := h.service.GetUserByEmail(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(user)
}

func (h *UserHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	var requestBody UpdateUserRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	id := c.Params("id")

	input := userservice.UpdateUserProfileInput{ID: id, Email: &requestBody.Email, Name: &requestBody.Name}
	err := h.service.UpdateUserProfile(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (h *UserHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	input := userservice.DeleteUserInput{ID: id}

	err := h.service.DeleteUser(input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(http.StatusOK)
}
