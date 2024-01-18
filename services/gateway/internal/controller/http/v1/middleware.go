package v1

import "github.com/gofiber/fiber/v2"

func (h *Handler) AuthMiddleware(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	_, err := h.auth.ParseToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	return c.Next()
}
