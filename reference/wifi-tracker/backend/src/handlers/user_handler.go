package handlers

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Service   service.UserServiceInterface
	Validator *validator.Validate
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		Service:   service,
		Validator: validator.New(),
	}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.Service.GetAll()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, users, "Successfully fetched users")
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID")
	}

	user, err := h.Service.GetByID(id)
	if err != nil {
		return utils.JSONError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.JSONSuccess(c, user, "Successfully fetched user")
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	dto := new(dtos.UserCreateDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Create(dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "User created successfully")
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID")
	}

	dto := new(dtos.UserUpdateDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.Validator.Struct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.Service.Update(id, dto); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "User updated successfully")
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, "Invalid UUID")
	}

	if err := h.Service.Delete(id); err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, nil, "User deleted successfully")
}
