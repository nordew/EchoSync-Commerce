package v1

import (
	"context"
	"gateway/internal/controller/http/dto"
	"github.com/gofiber/fiber/v2"

	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
)

func (h *Handler) signUp(c *fiber.Ctx) error {
	var input dto.SignUpRequest

	if err := c.BodyParser(&input); err != nil {
		writeInvalidJSONResponse(c)
		return err
	}

	_, err := h.grpcClient.SignUp(context.Background(), &nordew.SignUpRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		writeErrorResponse(c, "internal_error", err.Error())
		return err
	}

	c.Status(fiber.StatusCreated)

	return nil
}

func (h *Handler) signIn(c *fiber.Ctx) error {
	var input dto.SignInRequest

	if err := c.BodyParser(&input); err != nil {
		writeInvalidJSONResponse(c)
		return err
	}

	resp, err := h.grpcClient.SignIn(context.Background(), &nordew.SignInRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		writeErrorResponse(c, "internal_error", err.Error())
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})

	return nil
}
