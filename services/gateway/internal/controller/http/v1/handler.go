package v1

import (
	"gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
)

type Handler struct {
	logger logging.Logger

	grpcClient nordew.UserClient
}

func NewHandler(logger logging.Logger, grpcClient nordew.UserClient) *Handler {
	return &Handler{
		logger:     logger,
		grpcClient: grpcClient,
	}
}

func (h *Handler) Init() *fiber.App {
	app := fiber.New()

	auth := app.Group("/auth")
	auth.Post("/sign-up", h.signUp)
	auth.Get("/sign-in", h.signIn)

	return app
}

func writeInvalidJSONResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
		"error": "invalid json",
	})
}

func writeErrorResponse(c *fiber.Ctx, error, errorDescription string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{
		"error":             error,
		"error_description": errorDescription,
	})
}
