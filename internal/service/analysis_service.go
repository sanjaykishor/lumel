package service

import (
	"time"

	"github.com/sanjaykishor/lumel/internal/repository"
)

type AnalysisService struct {
	customerRepo *repository.CustomerRepository
}

func NewAnalysisService(customerRepo *repository.CustomerRepository) *AnalysisService {
	return &AnalysisService{
		customerRepo: customerRepo,
	}
}

type CustomerAnalysisParams struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (s *AnalysisService) GetCustomerAnalysis(params CustomerAnalysisParams) (*repository.CustomerAnalysisResult, error) {
	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		return nil, err
	}

	endDate = endDate.Add(24 * time.Hour)

	return s.customerRepo.GetCustomerAnalysis(startDate, endDate)
}
