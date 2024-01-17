package v1

import (
	"gateway/pkg/logging"
	"github.com/gofiber/fiber"
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

	return app
}

func writeInvalidJSONResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "invalid json",
	})
}

func writeErrorResponse(c *fiber.Ctx, error, errorDescription string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":             error,
		"error_description": errorDescription,
	})
}