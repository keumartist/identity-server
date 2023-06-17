package user

import (
	"net/http"

	userservice "art-sso/internal/service/user"

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

func (h *UserHandlerImpl) RegisterRoutes(app *fiber.App) {
	app.Post("/users", h.CreateUser)
	app.Get("/users/:id", h.GetUserByID)
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

func (h *UserHandlerImpl) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	input := userservice.GetUserByIDInput{ID: id}

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

	input := userservice.UpdateUserInput{ID: id, Email: &requestBody.Email, Password: &requestBody.Password}
	err := h.service.UpdateUser(input)
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
