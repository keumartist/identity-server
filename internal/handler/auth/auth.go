package auth

import (
	"net/http"

	errors "art-sso/internal/error"
	service "art-sso/internal/service/auth"

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
}

func (h *AuthHandlerImpl) SignUpWithEmail(c *fiber.Ctx) error {
	var requestBody SignUpWithEmailRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(errors.ErrInternal.Error())
	}

	tokens, err := h.authService.SignUpWithEmail(service.SignUpInput{Email: requestBody.Email, Password: requestBody.Password})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"idToken":      tokens.IdToken,
	})
}

func (h *AuthHandlerImpl) SignInWithEmail(c *fiber.Ctx) error {
	var requestBody SignInWithEmailRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	tokens, err := h.authService.SignInWithEmail(service.SignInInput{Email: requestBody.Email, Password: requestBody.Password})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"idToken":      tokens.IdToken,
	})
}
