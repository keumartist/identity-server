package handler

import (
	"net/http"

	domain "art-sso/internal/domain/user"
	service "art-sso/internal/service/user"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerImpl struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		service: service,
	}
}

func (h *UserHandlerImpl) RegisterRoutes(app *fiber.App) {
	app.Post("/users", h.CreateUserByEmail)
	app.Get("/users/:id", h.GetUserByID)
	app.Get("/users", h.GetUsers)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)
}

func (h *UserHandlerImpl) CreateUserByEmail(c *fiber.Ctx) error {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var requestBody RequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	createdUser, err := h.service.CreateUser(requestBody.Email, requestBody.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("User not found")
	}

	return c.Status(http.StatusCreated).JSON(createdUser)
}

func (h *UserHandlerImpl) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if user == nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}

	return c.JSON(user)
}

func (h *UserHandlerImpl) GetUsers(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(http.StatusBadRequest).SendString("Email query parameter is required")
	}

	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if user == nil {
		return c.Status(http.StatusNotFound).SendString("User not found")
	}

	return c.JSON(user)
}

func (h *UserHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	user := new(domain.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err := h.service.UpdateUser(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (h *UserHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = h.service.DeleteUser(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(http.StatusOK)
}
