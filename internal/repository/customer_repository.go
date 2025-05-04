package repository

import (
	"time"

	"github.com/sanjaykishor/lumel/internal/database"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

type CustomerAnalysisResult struct {
	TotalCustomers    int     `json:"total_customers"`
	TotalOrders       int     `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}

func (r *CustomerRepository) GetCustomerAnalysis(startDate, endDate time.Time) (*CustomerAnalysisResult, error) {
	result := &CustomerAnalysisResult{}

	if err := r.db.Model(&database.Order{}).
		Select("COUNT(DISTINCT customer_id) as total_customers").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&result.TotalCustomers).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&database.Order{}).
		Select("COUNT(*) as total_orders").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&result.TotalOrders).Error; err != nil {
		return nil, err
	}

	type AvgResult struct {
		Avg float64
	}

	var avgResult AvgResult

	if err := r.db.Model(&database.OrderItem{}).
		Select("COALESCE(AVG(order_items.unit_price * order_items.quantity * (1 - order_items.discount)), 0) as avg").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.date BETWEEN ? AND ?", startDate, endDate).
		Scan(&avgResult).Error; err != nil {
		return nil, err
	}

	result.AverageOrderValue = avgResult.Avg

	return result, nil
}
