package v1

import (
	"context"
	"gateway/internal/controller/http/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
)

func (h *Handler) createStore(c *fiber.Ctx) error {
	var input dto.CreateStoreRequest

	h.logger.Info("parsing body")
	if err := c.BodyParser(&input); err != nil {
		writeInvalidJSONResponse(c)
		return err
	}

	h.logger.Info("creating store", "name", input.Name)

	claims, _ := h.auth.ParseToken(c.Get("Authorization"))

	h.logger.Info("parsed claims", "sub", claims.Sub)

	parsedUUID, err := uuid.Parse(claims.Sub)
	if err != nil {
		writeErrorResponse(c, "internal_error", err.Error())
		return err
	}

	input.OwnerID = parsedUUID

	h.logger.Info("creating store", "name", input.Name, "owner_id", input.OwnerID)
	_, err = h.grpcStoreClient.CreateStore(context.Background(), &grpcStore.CreateStoreRequest{
		Name:    input.Name,
		OwnerId: input.OwnerID.String(),
	})
	if err != nil {
		writeErrorResponse(c, "internal_error", err.Error())
		return err
	}

	c.Status(fiber.StatusCreated)

	return nil
}
