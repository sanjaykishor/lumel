package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sanjaykishor/lumel/internal/config"
	"github.com/sanjaykishor/lumel/internal/database"
	"github.com/sanjaykishor/lumel/internal/utils"
)

func main() {
	csvPath := flag.String("csv", "", "Path to the CSV file")
	flag.Parse()

	cfg := config.Load()

	if *csvPath != "" {
		cfg.CSVPath = *csvPath
	}

	if _, err := os.Stat(cfg.CSVPath); os.IsNotExist(err) {
		log.Fatalf("CSV file does not exist: %s", cfg.CSVPath)
	}

	fmt.Printf("Starting data load from: %s\n", cfg.CSVPath)
	fmt.Printf("Database: %s@%s:%s/%s\n", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	startTime := time.Now()

	db, err := database.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	result, err := utils.ProcessCSVData(db, cfg.CSVPath)
	if err != nil {
		log.Fatalf("Error processing CSV data: %v", err)
	}

	duration := time.Since(startTime).Seconds()
	fmt.Println("\nData Load Results:")
	fmt.Println("==================")
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Start Time: %s\n", result.StartTime.Format(time.RFC3339))
	fmt.Printf("End Time: %s\n", result.EndTime.Format(time.RFC3339))
	fmt.Printf("Rows Processed: %d\n", result.RowsProcessed)
	fmt.Printf("Total Duration: %.2f seconds\n", duration)
}
