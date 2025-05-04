package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sanjaykishor/lumel/internal/api"
	"github.com/sanjaykishor/lumel/internal/config"
	"github.com/sanjaykishor/lumel/internal/database"
	"github.com/sanjaykishor/lumel/internal/repository"
	"github.com/sanjaykishor/lumel/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	customerRepo := repository.NewCustomerRepository(db)

	analysisService := service.NewAnalysisService(customerRepo)
	refreshService := service.NewRefreshService(db, cfg)

	handler := api.NewHandler(analysisService, refreshService)

	router := gin.Default()
	api.SetupRoutes(router, handler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
