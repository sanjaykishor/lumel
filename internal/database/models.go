package database

import (
	"time"
)

type Customer struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Category  string
	UnitPrice float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID            string `gorm:"primaryKey"`
	CustomerID    string
	Customer      Customer `gorm:"foreignKey:CustomerID"`
	Region        string
	Date          time.Time
	PaymentMethod string
	ShippingCost  float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   string
	Order     Order `gorm:"foreignKey:OrderID"`
	ProductID string
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int
	UnitPrice float64
	Discount  float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DataRefreshLog struct {
	ID            uint `gorm:"primaryKey"`
	StartTime     time.Time
	EndTime       time.Time
	Status        string
	Message       string
	RowsProcessed int
	CreatedAt     time.Time
}
