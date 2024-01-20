package v1

import (
	"context"
	"gateway/internal/controller/http/dto"
	"github.com/gofiber/fiber/v2"

	grpcUser "github.com/nordew/EchoSync-protos/gen/go/user"
)

func (h *Handler) signUp(c *fiber.Ctx) error {
	var input dto.SignUpRequest

	if err := c.BodyParser(&input); err != nil {
		writeInvalidJSONResponse(c)
		return err
	}

	_, err := h.grpcUserClient.SignUp(context.Background(), &grpcUser.SignUpRequest{
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

	resp, err := h.grpcUserClient.SignIn(context.Background(), &grpcUser.SignInRequest{
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

func (h *Handler) refresh(c *fiber.Ctx) error {
	refreshToken := c.Get("refresh_token")

	resp, err := h.grpcUserClient.Refresh(context.Background(), &grpcUser.RefreshRequest{
		RefreshToken: refreshToken,
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
