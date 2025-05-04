package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjaykishor/lumel/internal/service"
)

type CustomerAnalysisQuery struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type Handler struct {
	analysisService *service.AnalysisService
	refreshService  *service.RefreshService
}

func NewHandler(
	analysisService *service.AnalysisService,
	refreshService *service.RefreshService,
) *Handler {
	return &Handler{
		analysisService: analysisService,
		refreshService:  refreshService,
	}
}

func (h *Handler) GetCustomerAnalysis(c *gin.Context) {
	var query CustomerAnalysisQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters: " + err.Error()})
		return
	}

	_, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start_date format. Use YYYY-MM-DD",
		})
		return
	}

	_, err = time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end_date format. Use YYYY-MM-DD",
		})
		return
	}

	result, err := h.analysisService.GetCustomerAnalysis(service.CustomerAnalysisParams{
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analysis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) RefreshData(c *gin.Context) {
	result, err := h.refreshService.RefreshData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetRefreshHistory(c *gin.Context) {
	logs, err := h.refreshService.GetRefreshHistory(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get refresh history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
