package user

import (
	"errors"
	"net/http"

	customerror "art-sso/internal/error"
	tokenservice "art-sso/internal/service/token"
	userservice "art-sso/internal/service/user"

	middleware "art-sso/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerImpl struct {
	service userservice.UserService
}

func NewUserHandler(service userservice.UserService) UserHandler {
	return &UserHandlerImpl{
		service: service,
	}
}

func (h *UserHandlerImpl) RegisterRoutes(app *fiber.App, tokenService tokenservice.TokenService) {
	app.Get("/api/v1/users/me", middleware.TokenValidationMiddleware(tokenService), h.GetMe)
	app.Get("/api/v1/users", h.GetUsers)
	app.Put("/api/v1/users/me", middleware.TokenValidationMiddleware(tokenService), h.UpdateMeUserProfile)
	app.Delete("/api/v1/users/me", middleware.TokenValidationMiddleware(tokenService), h.DeleteMeUser)
}

func (h *UserHandlerImpl) GetMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(customerror.ErrUnauthorized)
	}

	input := userservice.GetUserByIDInput{ID: userId}

	user, err := h.service.GetUserByID(input)
	if err != nil {
		if errors.Is(err, customerror.ErrUserNotFound) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrUserNotFound)
		}
		return c.Status(http.StatusInternalServerError).JSON(customerror.ErrInternal)
	}

	return c.JSON(user)
}

func (h *UserHandlerImpl) GetUsers(c *fiber.Ctx) error {
	email := c.Query("email")
	input := userservice.GetUserByEmailInput{Email: email}

	user, err := h.service.GetUserByEmail(input)
	if err != nil {
		if errors.Is(err, customerror.ErrUserNotFound) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrUserNotFound)
		}
		return c.Status(http.StatusInternalServerError).JSON(customerror.ErrInternal)
	}

	return c.JSON(user)
}

func (h *UserHandlerImpl) UpdateMeUserProfile(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(customerror.ErrUnauthorized)
	}

	var requestBody UpdateUserProfileRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
	}

	input := userservice.UpdateUserProfileInput{ID: userId, Name: &requestBody.Name}
	err := h.service.UpdateUserProfile(input)
	if err != nil {
		if errors.Is(err, customerror.ErrBadRequest) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
		}
		return c.Status(http.StatusInternalServerError).JSON(customerror.ErrInternal)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *UserHandlerImpl) DeleteMeUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).SendString("Not authorized")
	}

	input := userservice.DeleteUserInput{ID: userId}

	err := h.service.DeleteUser(input)
	if err != nil {
		if errors.Is(err, customerror.ErrBadRequest) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
		}
		return c.Status(http.StatusInternalServerError).JSON(customerror.ErrInternal)
	}

	return c.SendStatus(http.StatusOK)
}
