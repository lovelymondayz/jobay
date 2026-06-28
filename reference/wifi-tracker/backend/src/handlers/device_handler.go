package handlers

import (
	"regexp"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DevicesHandler struct {
	Service   service.DevicesServiceInterface
	Validator *validator.Validate
}

func NewDevicesHandler(service service.DevicesServiceInterface) *DevicesHandler {
	v := validator.New()

	// Custom validator untuk MAC address
	v.RegisterValidation("mac", func(fl validator.FieldLevel) bool {
		mac := fl.Field().String()
		matched, _ := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, mac)
		return matched
	})

	return &DevicesHandler{
		Service:   service,
		Validator: v,
	}
}

func (h *DevicesHandler) GetAll(c *fiber.Ctx) error {
	devices, err := h.Service.GetAll()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, devices, "Success to find all devices")
}

func (h *DevicesHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID format")
	}

	device, err := h.Service.GetByID(id)
	if err != nil {
		return utils.JSONError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.JSONSuccess(c, device, "Success to find device")
}

func (h *DevicesHandler) Create(c *fiber.Ctx) error {
	dto := new(dtos.DeviceDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Create(dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Device created successfully")
}

func (h *DevicesHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID format")
	}

	dto := new(dtos.DeviceDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Update(id, dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Device updated successfully")
}

func (h *DevicesHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID format")
	}

	if err := h.Service.Delete(id); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Device deleted successfully")
}

