package middleware

import (
	"log"
	"runtime/debug"
	"wifi-tracker-be/src/utils"

	"github.com/gofiber/fiber/v2"
)

func CustomRecoverPanic() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log panic dan stack trace
				stackTrace := debug.Stack()
				log.Printf("[PANIC] %v\n%s\n", r, stackTrace)

				// Kirim respons error kustom
				_ = utils.JSONError(c, fiber.StatusInternalServerError, "Internal Server Error")
			}
		}()

		// Lanjutkan ke handler berikutnya
		return c.Next()
	}
}
