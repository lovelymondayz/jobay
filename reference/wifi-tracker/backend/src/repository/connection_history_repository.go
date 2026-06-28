package repository

import (
	"math"
	"time"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConnectionHistoriesRepositoryInterface interface {
	FindAll(tx *gorm.DB, search string, dateFrom string, dateTo string, page int, size int) (datas []dtos.ConnectionHistoriesResponse, totals int, maxPages int, err error) 
	FindById(tx *gorm.DB, id uuid.UUID) (*models.ConnectionHistories, error)
	Create(tx *gorm.DB, connectionHistory *models.ConnectionHistories) error
	Update(tx *gorm.DB, connectionHistory *models.ConnectionHistories) error
	Delete(tx *gorm.DB, id uuid.UUID) error
	ExistsToday(tx *gorm.DB, deviceID uuid.UUID, date time.Time) (*dtos.ConnectionHistoriesDataExist, error)
}

type ConnectionHistoriesRepository struct{}

func NewConnectionHistoriesRepository() ConnectionHistoriesRepositoryInterface {
	return &ConnectionHistoriesRepository{}
}

func (r *ConnectionHistoriesRepository) FindAll(
	tx *gorm.DB,
	search string,
	dateFrom string,
	dateTo string,
	page int,
	size int,
) (datas []dtos.ConnectionHistoriesResponse, totals int, maxPages int, err error) {
	var connectionHistories []dtos.ConnectionHistoriesResponse
	var total int64 // gunakan int64 untuk Count()

	query := tx.
		Table("connection_histories as c").
		Select(`
			c.connection_history_id,
			u.name as user_name,
			d.name as device_name,
			a.name as to_aps,
			b.name as from_aps,
			c.created_at
		`).
		Joins("LEFT JOIN users as u ON c.user_id = u.user_id").
		Joins("LEFT JOIN devices as d ON c.device_id = d.device_id").
		Joins("LEFT JOIN access_points as a ON c.to_aps = a.access_point_id").
		Joins("LEFT JOIN access_points as b ON c.from_aps = b.access_point_id")

	// Filter search
	if search != "" {
		query = query.Where(`(
			u.name ILIKE ? OR
			d.name ILIKE ? OR
			a.name ILIKE ? OR
			b.name ILIKE ?
		)
		`, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Filter tanggal
	if dateFrom != "" && dateTo != "" {
		query = query.Where("c.created_at BETWEEN ? AND ?", dateFrom+"+07", dateTo+"+07")
	}

	// Hitung total sebelum pagination
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, 0, err
	}

	// Pagination & eksekusi query
	err = query.
		Order("c.created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&connectionHistories).Error
	if err != nil {
		return nil, 0, 0, err
	}

	// Hitung maxPages
	maxPages = int(math.Ceil(float64(total) / float64(size)))

	return connectionHistories, int(total), maxPages, nil
}

func (r *ConnectionHistoriesRepository) FindById(tx *gorm.DB, id uuid.UUID) (*models.ConnectionHistories, error) {
	var connectionHistory models.ConnectionHistories
	if err := tx.Where("connection_history_id = ?", id).First(&connectionHistory).Error; err != nil {
		return nil, err
	}
	return &connectionHistory, nil
}

func (r *ConnectionHistoriesRepository) Create(tx *gorm.DB, connectionHistory *models.ConnectionHistories) error {
	if err := tx.Create(connectionHistory).Error; err != nil {
		return err
	}
	return nil
}

func (r *ConnectionHistoriesRepository) Update(tx *gorm.DB, connectionHistory *models.ConnectionHistories) error {
	if err := tx.Save(connectionHistory).Error; err != nil {
		return err
	}
	return nil
}

func (r *ConnectionHistoriesRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if err := tx.Where("connection_history_id = ?", id).Delete(&models.ConnectionHistories{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ConnectionHistoriesRepository) ExistsToday(tx *gorm.DB, deviceID uuid.UUID, date time.Time) (*dtos.ConnectionHistoriesDataExist, error) {
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var datas dtos.ConnectionHistoriesDataExist

	err := tx.Table("connection_histories AS c").
		Where("c.device_id = ? AND c.created_at >= ? AND c.created_at < ?", deviceID, startOfDay, endOfDay).
		Select(`
			c.connection_history_id,
			c.to_aps as to_aps,
			d.mac_address as to_mac_address_device
		`).
		Joins("LEFT JOIN access_points as a ON c.to_aps = a.access_point_id").
		Joins("LEFT JOIN devices as d ON c.device_id = d.device_id").
		Order("c.created_at DESC").
		Limit(1).
		Scan(&datas).Error

	return &datas, err
}
