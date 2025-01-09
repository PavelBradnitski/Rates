package services

import (
	"context"

	"github.com/PavelBradnitski/Rates/pkg/models"
	"github.com/PavelBradnitski/Rates/pkg/repositories"
)

type RateService struct {
	repo *repositories.RateRepository
}

func NewRateService(repo *repositories.RateRepository) *RateService {
	return &RateService{repo: repo}
}

func (s *RateService) GetAllRates(ctx context.Context) ([]models.Rate, error) {
	return s.repo.GetAllRates(ctx)
}

func (s *RateService) GetRateByDate(ctx context.Context, date string) ([]models.Rate, error) {
	return s.repo.GetRateByDate(ctx, date)
}
