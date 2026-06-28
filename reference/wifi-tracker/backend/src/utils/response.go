package utils

import (
	"wifi-tracker-be/src/dtos"

	"github.com/gofiber/fiber/v2"
)

func JSONSuccess(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(dtos.WebResponseGeneral{
		Status:  true,
		Message: message,
		Data:    data,
	})
	
}

// ValidationErrorResponse menampilkan response JSON standar untuk error validasi
func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	message := FormatValidationError(err)
	return JSONError(c, fiber.StatusBadRequest, message)
}

// Untuk response error
func JSONError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(dtos.WebResponseGeneral{
		Status:  false,
		Message: message,
		Data:    nil,
	})
}



// return c.Status(fiber.StatusOK).JSON(models.ActiveDevice{
// 	ID: id,
// 	MACAddress: mac_address,
// 	APName: ap_name,        
// 	TxBytes:tx_bytes,        
// 	RxBytes:rx_bytes,
// 	ConnectedAt:connected_at,
// 	DisconnectedAt:disconnected_at,
// 	IsActive:false,
// })