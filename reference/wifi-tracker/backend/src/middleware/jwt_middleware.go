package middleware

import (
	"strings"
	"wifi-tracker-be/src/utils"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return utils.JSONError(c, fiber.StatusUnauthorized, "Missing token")
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return utils.JSONError(c, fiber.StatusUnauthorized, "Invalid token format")
        }

        userID, err := utils.ParseJWT(parts[1])
        if err != nil {
            return utils.JSONError(c, fiber.StatusUnauthorized, "Invalid token")
        }

        c.Locals("userID", userID)
        return c.Next()
    }
}
