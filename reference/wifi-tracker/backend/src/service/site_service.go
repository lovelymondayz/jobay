package service

import (
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SiteServiceInterface interface {
	FindAllSites() ([]models.Sites, error)
	FindSiteByID(siteID string) (*models.Sites, error)
	CreateSite(site models.Sites) error
	UpdateSite(site models.Sites) error
	DeleteSite(siteID string) error
}

type SiteService struct {
	repo repository.SiteRepositoryInterface
	DB   *gorm.DB
}

func NewSiteService(repo repository.SiteRepositoryInterface, db *gorm.DB) SiteServiceInterface {
	return &SiteService{
		repo: repo,
		DB:   db,
	}
}

func (s *SiteService) FindAllSites() ([]models.Sites, error) {
	return s.repo.FindAll(s.DB)
}

func (s *SiteService) FindSiteByID(siteID string) (*models.Sites, error) {
	parseId, err := uuid.Parse(siteID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindById(s.DB, parseId)
}

func (s *SiteService) CreateSite(site models.Sites) error {
	return s.repo.Create(s.DB, &site)
}

func (s *SiteService) UpdateSite(site models.Sites) error {
	return s.repo.Update(s.DB, &site)
}

func (s *SiteService) DeleteSite(siteID string) error {
	parseId, err := uuid.Parse(siteID)
	if err != nil {
		return err
	}
	return s.repo.Delete(s.DB, parseId)
}
