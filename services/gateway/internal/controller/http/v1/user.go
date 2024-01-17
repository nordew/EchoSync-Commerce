package v1

import (
	"gateway/internal/controller/http/dto"

	"github.com/gofiber/fiber"
	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
)

func (h *Handler) signUp(c *fiber.Ctx) {
	var input dto.SignUpRequest

	if err := c.BodyParser(&input); err != nil {
		writeInvalidJSONResponse(c)
	}

	_, err := h.grpcClient.SignUp(c.Context(), &nordew.SignUpRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		writeErrorResponse(c, "internal_error", err.Error())
	}

	c.Status(fiber.StatusCreated)
}
