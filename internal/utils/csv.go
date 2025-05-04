package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sanjaykishor/lumel/internal/database"
	"gorm.io/gorm"
)

type CSVRow struct {
	OrderID         string
	ProductID       string
	CustomerID      string
	ProductName     string
	Category        string
	Region          string
	DateOfSale      string
	QuantitySold    string
	UnitPrice       string
	Discount        string
	ShippingCost    string
	PaymentMethod   string
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string
}

func ParseCSV(filePath string) ([]CSVRow, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV content: %w", err)
	}

	if len(records) < 2 { 
		return nil, errors.New("CSV file has insufficient data")
	}

	records = records[1:]

	rows := make([]CSVRow, len(records))
	for i, record := range records {
		if len(record) < 15 {
			return nil, fmt.Errorf("row %d has insufficient columns", i+1)
		}

		rows[i] = CSVRow{
			OrderID:         record[0],
			ProductID:       record[1],
			CustomerID:      record[2],
			ProductName:     record[3],
			Category:        record[4],
			Region:          record[5],
			DateOfSale:      record[6],
			QuantitySold:    record[7],
			UnitPrice:       record[8],
			Discount:        record[9],
			ShippingCost:    record[10],
			PaymentMethod:   record[11],
			CustomerName:    record[12],
			CustomerEmail:   record[13],
			CustomerAddress: record[14],
		}
	}

	return rows, nil
}

func ProcessCSVData(db *gorm.DB, filePath string) (*database.DataRefreshLog, error) {
	refreshLog := &database.DataRefreshLog{
		StartTime: time.Now(),
		Status:    "PROCESSING",
	}

	if err := db.Create(refreshLog).Error; err != nil {
		return nil, fmt.Errorf("failed to create refresh log: %w", err)
	}

	rows, err := ParseCSV(filePath)
	if err != nil {
		refreshLog.Status = "FAILED"
		refreshLog.Message = fmt.Sprintf("Error parsing CSV: %v", err)
		refreshLog.EndTime = time.Now()
		db.Save(refreshLog)
		return refreshLog, err
	}

	tx := db.Begin()

	customers := make(map[string]bool)
	products := make(map[string]bool)
	orders := make(map[string]bool)

	for _, row := range rows {
		if !customers[row.CustomerID] {
			customer := database.Customer{
				ID:      row.CustomerID,
				Name:    row.CustomerName,
				Email:   row.CustomerEmail,
				Address: row.CustomerAddress,
			}

			if err := tx.Where("id = ?", customer.ID).FirstOrCreate(&customer).Error; err != nil {
				tx.Rollback()
				refreshLog.Status = "FAILED"
				refreshLog.Message = fmt.Sprintf("Error creating customer: %v", err)
				refreshLog.EndTime = time.Now()
				db.Save(refreshLog)
				return refreshLog, err
			}
			customers[row.CustomerID] = true
		}

		if !products[row.ProductID] {
			unitPrice, _ := strconv.ParseFloat(row.UnitPrice, 64)
			product := database.Product{
				ID:        row.ProductID,
				Name:      row.ProductName,
				Category:  row.Category,
				UnitPrice: unitPrice,
			}

			if err := tx.Where("id = ?", product.ID).FirstOrCreate(&product).Error; err != nil {
				tx.Rollback()
				refreshLog.Status = "FAILED"
				refreshLog.Message = fmt.Sprintf("Error creating product: %v", err)
				refreshLog.EndTime = time.Now()
				db.Save(refreshLog)
				return refreshLog, err
			}
			products[row.ProductID] = true
		}

		if !orders[row.OrderID] {
			date, _ := time.Parse("2006-01-02", row.DateOfSale)
			shippingCost, _ := strconv.ParseFloat(row.ShippingCost, 64)
			order := database.Order{
				ID:            row.OrderID,
				CustomerID:    row.CustomerID,
				Region:        row.Region,
				Date:          date,
				PaymentMethod: row.PaymentMethod,
				ShippingCost:  shippingCost,
			}

			if err := tx.Where("id = ?", order.ID).FirstOrCreate(&order).Error; err != nil {
				tx.Rollback()
				refreshLog.Status = "FAILED"
				refreshLog.Message = fmt.Sprintf("Error creating order: %v", err)
				refreshLog.EndTime = time.Now()
				db.Save(refreshLog)
				return refreshLog, err
			}
			orders[row.OrderID] = true
		}

		quantity, _ := strconv.Atoi(row.QuantitySold)
		unitPrice, _ := strconv.ParseFloat(row.UnitPrice, 64)
		discount, _ := strconv.ParseFloat(row.Discount, 64)

		orderItem := database.OrderItem{
			OrderID:   row.OrderID,
			ProductID: row.ProductID,
			Quantity:  quantity,
			UnitPrice: unitPrice,
			Discount:  discount,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			refreshLog.Status = "FAILED"
			refreshLog.Message = fmt.Sprintf("Error creating order item: %v", err)
			refreshLog.EndTime = time.Now()
			db.Save(refreshLog)
			return refreshLog, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		refreshLog.Status = "FAILED"
		refreshLog.Message = fmt.Sprintf("Error committing transaction: %v", err)
		refreshLog.EndTime = time.Now()
		db.Save(refreshLog)
		return refreshLog, err
	}

	refreshLog.Status = "COMPLETED"
	refreshLog.Message = "Data loaded successfully"
	refreshLog.RowsProcessed = len(rows)
	refreshLog.EndTime = time.Now()
	db.Save(refreshLog)

	return refreshLog, nil
}
