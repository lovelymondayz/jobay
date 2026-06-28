package handlers

import (
	"strconv"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RoleHandlerInterface interface {
	FindAll(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type RoleHandler struct {
	Service   service.RoleServiceInterface
	Validator *validator.Validate
}

func NewRoleHandler(service service.RoleServiceInterface) RoleHandlerInterface {
	return &RoleHandler{
		Service:   service,
		Validator: validator.New(),
	}
}
func (h *RoleHandler) FindAll(c *fiber.Ctx) error {
	roles, err := h.Service.GetAll()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, roles, "Success to find all data")
}

func (h *RoleHandler) FindByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	role, err := h.Service.GetByID(uint(id))
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, role, "Success to find data")
}

func (h *RoleHandler) Create(c *fiber.Ctx) error {
	dto := new(dtos.RoleDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.Validator.Struct(dto); err != nil {
		errorMsg := utils.FormatValidationError(err)
		return utils.JSONError(c, fiber.StatusBadRequest, errorMsg)
	}

	if err := h.Service.Create(dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Success to create data")
}

func (h *RoleHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid ID")
	}

	dto := new(dtos.RoleDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.Validator.Struct(dto); err != nil {
		errorMsg := utils.FormatValidationError(err)
		return utils.JSONError(c, fiber.StatusBadRequest, errorMsg)
	}

	if err := h.Service.Update(uint(id), dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Success to update data")
}

func (h *RoleHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid ID")
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "Success to delete data")
}
