package handlers

import (
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/gofiber/fiber/v2"
)

type UnifiIntegrationHandlerInterface interface {
	GetClients(c *fiber.Ctx) error
	GetClients2(c *fiber.Ctx) error
	GetDevices(c *fiber.Ctx) error
	GetDevices2(c *fiber.Ctx) error
	GetClientByMACAP(c *fiber.Ctx) error
	FindAllHistoryConnection(c *fiber.Ctx) error
}

type UnifiIntegrationHandler struct {
	unifiService       service.UnifiIntegrationServiceInterface
	anotherUnfiService service.UnifiIntegrationServiceWoUnifiPollerInterface
}

func NewUnifiIntegrationHandler(
	unifiService service.UnifiIntegrationServiceInterface,
	anotherUnfiService service.UnifiIntegrationServiceWoUnifiPollerInterface,
) UnifiIntegrationHandlerInterface {
	return &UnifiIntegrationHandler{
		unifiService:       unifiService,
		anotherUnfiService: anotherUnfiService,
	}
}

func (h *UnifiIntegrationHandler) GetClients(c *fiber.Ctx) error {
	datas, err := h.unifiService.GetClients()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, datas, "Success to get active clients")
}

func (h *UnifiIntegrationHandler) GetClients2(c *fiber.Ctx) error {
	datas, err := h.anotherUnfiService.GetClientsWithNearbyAPs()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, datas, "Success to get active clients")
}

func (h *UnifiIntegrationHandler) GetDevices(c *fiber.Ctx) error {
	datas, err := h.unifiService.GetDevices()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, datas, "Success to get Devices Active")
}
func (h *UnifiIntegrationHandler) GetDevices2(c *fiber.Ctx) error {
	datas, err := h.anotherUnfiService.GetDevices()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, datas, "Success to get Devices Active")
}

func (h *UnifiIntegrationHandler) GetClientByMACAP(c *fiber.Ctx) error {
	macAP := c.Params("macAP")

	datas, err := h.anotherUnfiService.GetClientByMACAP(macAP)
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, datas, "Success to get active clients By MAC AP")
}

func (h *UnifiIntegrationHandler) FindAllHistoryConnection(c *fiber.Ctx) error {
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page := c.QueryInt("page")
	size := c.QueryInt("size")

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 10
	}
	datas, err := h.anotherUnfiService.FindAllHistoryConnection(search, startDate, endDate, page, size)
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, datas, "Success to get history connection")
}
