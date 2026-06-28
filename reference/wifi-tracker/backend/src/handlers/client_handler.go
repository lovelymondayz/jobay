package handlers

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClientsHandler struct {
	Service   service.ClientServiceInterface
	Validator *validator.Validate
}

func NewClientsHandler(service service.ClientServiceInterface) *ClientsHandler {
	return &ClientsHandler{
		Service:   service,
		Validator: validator.New(),
	}
}

func (h *ClientsHandler) GetAll(c *fiber.Ctx) error {
	clients, err := h.Service.GetAll()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, clients, "Success to find all clients")
}

func (h *ClientsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID format")
	}

	client, err := h.Service.GetByID(id)
	if err != nil {
		return utils.JSONError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.JSONSuccess(c, client, "Success to find client")
}

func (h *ClientsHandler) Create(c *fiber.Ctx) error {
	dto := new(dtos.ClientDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Create(dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Client created successfully")
}

func (h *ClientsHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID format")
	}

	dto := new(dtos.ClientDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Update(id, dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Client updated successfully")
}

func (h *ClientsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID format")
	}

	if err := h.Service.Delete(id); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Client deleted successfully")
}
