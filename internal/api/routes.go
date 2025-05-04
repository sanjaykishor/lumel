package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handler.HealthCheck)

		v1.GET("/analysis/customer", handler.GetCustomerAnalysis)

		v1.POST("/data/refresh", handler.RefreshData)
		v1.GET("/data/refresh/history", handler.GetRefreshHistory)
	}
}
