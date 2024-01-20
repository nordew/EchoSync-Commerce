package v1

import (
	"gateway/pkg/auth"
	"gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"

	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
	grpcUser "github.com/nordew/EchoSync-protos/gen/go/user"
)

type Handler struct {
	logger logging.Logger

	grpcUserClient    grpcUser.UserClient
	grpcStoreClient   grpcStore.StoreServiceClient
	grpcProductClient grpcStore.ProductServiceClient

	auth auth.Authenticator
}

func NewHandler(logger logging.Logger, grpcUserClient grpcUser.UserClient, grpcStoreClient grpcStore.StoreServiceClient, grpcProductClient grpcStore.ProductServiceClient, auth auth.Authenticator) *Handler {
	return &Handler{
		logger:            logger,
		grpcUserClient:    grpcUserClient,
		grpcStoreClient:   grpcStoreClient,
		grpcProductClient: grpcProductClient,
		auth:              auth,
	}
}

func (h *Handler) Init() *fiber.App {
	app := fiber.New()

	auth := app.Group("/auth")
	{
		auth.Post("/sign-up", h.signUp)
		auth.Get("/sign-in", h.signIn)
		auth.Get("/refresh", h.refresh)
	}

	market := app.Group("/market")
	market.Use(h.AuthMiddleware)
	{
		market.Post("/store", h.createStore)
	}

	product := app.Group("/product")
	product.Use(h.AuthMiddleware)
	{
		product.Post("/create", h.createProduct)
	}

	return app
}

func writeInvalidJSONResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Invalid JSON",
	})
}

func writeErrorResponse(c *fiber.Ctx, errType, errDescription string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":             errType,
		"error_description": errDescription,
	})
}
