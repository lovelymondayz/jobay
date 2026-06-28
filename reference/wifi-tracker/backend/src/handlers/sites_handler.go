package handlers

import (
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SitesHandlerInterface interface {
	FindAllSites(c *fiber.Ctx) error
	FindSiteById(c *fiber.Ctx) error
	CreateSite(c *fiber.Ctx) error
	UpdateSite(c *fiber.Ctx) error
	DeleteSite(c *fiber.Ctx) error
}

type SitesHandler struct {
	siteService service.SiteServiceInterface
	validate *validator.Validate
}

func NewSitesHandler(siteService service.SiteServiceInterface) SitesHandlerInterface {
	return &SitesHandler{
		siteService: siteService,
		validate:    validator.New(),
	}
}

func (h SitesHandler) FindAllSites(c *fiber.Ctx) error {
	datas, err := h.siteService.FindAllSites()
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, datas, "Success to find all data")
}

func (h SitesHandler) FindSiteById(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	datas, err := h.siteService.FindSiteByID(siteID)
	if err != nil {
		return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, datas, "Success to find data")
}

func (h SitesHandler) CreateSite(c *fiber.Ctx) error {
	var site models.Sites
	if err := c.BodyParser(&site); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(&site); err != nil {
		errorMsg := utils.FormatValidationError(err)
		return utils.JSONError(c, fiber.StatusBadRequest, errorMsg)
	}
	return utils.JSONSuccess(c, h.siteService.CreateSite(site), "Success to create data")
}

func (h SitesHandler) UpdateSite(c *fiber.Ctx) error {
	var site models.Sites
	if err := c.BodyParser(&site); err != nil {
		return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.JSONSuccess(c, h.siteService.UpdateSite(site), "Success to update data")
}

func (h SitesHandler) DeleteSite(c *fiber.Ctx) error {
	siteID := c.Params("siteID")
	return utils.JSONSuccess(c, h.siteService.DeleteSite(siteID), "Success to delete data")
}