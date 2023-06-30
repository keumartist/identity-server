package auth

import (
	"net/http"

	customerror "art-sso/internal/error"
	service "art-sso/internal/service/auth"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AuthHandlerImpl struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authService: authService,
	}
}

func (h *AuthHandlerImpl) RegisterRoutes(app *fiber.App) {
	app.Post("/signup", h.SignUpWithEmail)
	app.Post("/signin", h.SignInWithEmail)
	app.Post("/verification", h.VerifyEmail)
}

func (h *AuthHandlerImpl) SignUpWithEmail(c *fiber.Ctx) error {
	var requestBody SignUpWithEmailRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
	}

	message, err := h.authService.SignUpWithEmail(service.SignUpInput{Email: requestBody.Email, Password: requestBody.Password})
	if err != nil {
		if errors.Is(err, customerror.ErrEmailInUse) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrEmailInUse)
		}
		return c.Status(http.StatusInternalServerError).JSON(customerror.ErrInternal)
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
}

func (h *AuthHandlerImpl) SignInWithEmail(c *fiber.Ctx) error {
	var requestBody SignInWithEmailRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
	}

	tokens, err := h.authService.SignInWithEmail(service.SignInInput{Email: requestBody.Email, Password: requestBody.Password})
	if err != nil {
		if errors.Is(err, customerror.ErrInvalidCredentials) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrInvalidCredentials)
		}
		return c.Status(http.StatusInternalServerError).JSON(customerror.ErrInternal)
	}

	return c.JSON(fiber.Map{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"idToken":      tokens.IdToken,
	})
}

func (h *AuthHandlerImpl) SignInWithGoogle(c *fiber.Ctx) error {
	var requestBody SignInWithGoogleRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
	}

	tokens, err := h.authService.SignInWithGoogle(service.SignInWithGoogleInput{Code: requestBody.Code})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error()) // TODO: Fix error handling
	}

	return c.JSON(fiber.Map{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"idToken":      tokens.IdToken,
	})
}

func (h *AuthHandlerImpl) VerifyEmail(c *fiber.Ctx) error {
	var requestBody VerifyEmailRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
	}

	err := h.authService.VerifyEmailCode(service.VerifyEmailCodeInput{Email: requestBody.Email, Code: requestBody.Code})
	if err != nil {
		if errors.Is(err, customerror.ErrInvalidVerificationCode) {
			return c.Status(http.StatusBadRequest).JSON(customerror.ErrInvalidVerificationCode)
		}
		return c.Status(http.StatusBadRequest).JSON(customerror.ErrBadRequest)
	}

	return c.Status(http.StatusNoContent).SendString("")
}
