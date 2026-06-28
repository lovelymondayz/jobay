package handlers

import (
	"wifi-tracker-be/src/utils"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type healthHandler struct {
}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

func (h *healthHandler) HealthCheck(c *fiber.Ctx) error {
	return utils.JSONSuccess(c, nil, "Success to check health")
}
