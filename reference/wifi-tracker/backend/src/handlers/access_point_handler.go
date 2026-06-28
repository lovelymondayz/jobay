package handlers

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccessPointHandler struct {
	Service   service.AccessPointServiceInterface
	Validator *validator.Validate
}

func NewAccessPointHandler(service service.AccessPointServiceInterface) *AccessPointHandler {
	return &AccessPointHandler{
		Service:   service,
		Validator: validator.New(),
	}
}

func (h *AccessPointHandler) GetAll(c *fiber.Ctx) error {
	aps, err := h.Service.GetAll()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, aps, "Access points fetched successfully")
}

func (h *AccessPointHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID")
	}

	ap, err := h.Service.GetByID(id)
	if err != nil {
		return utils.JSONError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.JSONSuccess(c, ap, "Access point fetched successfully")
}

func (h *AccessPointHandler) Create(c *fiber.Ctx) error {
	dto := new(dtos.AccessPointDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Create(dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Access point created successfully")
}

func (h *AccessPointHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID")
	}

	dto := new(dtos.AccessPointDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Update(id, dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Access point updated successfully")
}

func (h *AccessPointHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID")
	}

	if err := h.Service.Delete(id); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Access point deleted successfully")
}
