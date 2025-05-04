package service

import (
	"time"

	"github.com/sanjaykishor/lumel/internal/config"
	"github.com/sanjaykishor/lumel/internal/database"
	"github.com/sanjaykishor/lumel/internal/utils"
	"gorm.io/gorm"
)

type RefreshService struct {
	db     *gorm.DB
	config *config.Config
}

func NewRefreshService(db *gorm.DB, config *config.Config) *RefreshService {
	return &RefreshService{
		db:     db,
		config: config,
	}
}

type RefreshResult struct {
	Success       bool      `json:"success"`
	Message       string    `json:"message"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	RowsProcessed int       `json:"rows_processed"`
}

func (s *RefreshService) RefreshData() (*RefreshResult, error) {
	refreshLog, err := utils.ProcessCSVData(s.db, s.config.CSVPath)
	if err != nil {
		return &RefreshResult{
			Success:   false,
			Message:   err.Error(),
			StartTime: refreshLog.StartTime,
			EndTime:   refreshLog.EndTime,
		}, err
	}

	return &RefreshResult{
		Success:       true,
		Message:       refreshLog.Message,
		StartTime:     refreshLog.StartTime,
		EndTime:       refreshLog.EndTime,
		RowsProcessed: refreshLog.RowsProcessed,
	}, nil
}

func (s *RefreshService) GetRefreshHistory(limit int) ([]database.DataRefreshLog, error) {
	var logs []database.DataRefreshLog

	if err := s.db.Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
