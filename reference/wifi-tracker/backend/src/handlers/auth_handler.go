package handlers

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerInterface interface {
	// Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService     service.AuthServiceInterface
	validate *validator.Validate
}

// Constructor
func NewAuthHandler(authService service.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// func (h *authHandler) Register(c *fiber.Ctx) error {
// 	var input models.User

// 	if err := c.BodyParser(&input); err != nil {
// 		return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
// 	}

// 	// Generate UUID
// 	input.ID = uuid.New()

// 	// Hash password
// 	hashedPassword, err := utils.HashPassword(input.Password)
// 	if err != nil {
// 		return utils.JSONError(c, fiber.StatusInternalServerError, "Password hashing failed")
// 	}
// 	input.Password = hashedPassword

// 	// Save to DB
// 	if err := h.repo.Create(&input); err != nil {
// 		return utils.JSONError(c, fiber.StatusInternalServerError, "Could not create user")
// 	}

// 	return utils.JSONSuccess(c, nil, "Success to create data")
// }

func (h *authHandler) Login(c *fiber.Ctx) error {
	var data dtos.LoginRequest
	if err := c.BodyParser(&data); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(&data); err != nil {
		errorMsg := utils.FormatValidationError(err)
		return utils.JSONError(c, fiber.StatusBadRequest, errorMsg)
	}

	token, err := h.authService.Login(&data)
	if err != nil {
		return utils.JSONError(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.JSONSuccess(c, fiber.Map{"token": token}, "Success to login")
}
