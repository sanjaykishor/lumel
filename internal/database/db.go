package database

import (
	"fmt"

	"github.com/sanjaykishor/lumel/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDBConnection creates a new database connection
func NewDBConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&Customer{}, &Product{}, &Order{}, &OrderItem{}, &DataRefreshLog{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
