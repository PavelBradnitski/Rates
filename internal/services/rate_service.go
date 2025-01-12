package services

import (
	"context"

	"github.com/PavelBradnitski/Rates/internal/models"
	"github.com/PavelBradnitski/Rates/internal/repositories"
)

type RateService struct {
	repo repositories.RateRepositoryInterface
}

type RateServiceInterface interface {
	GetAllRates(ctx context.Context) ([]models.Rate, error)
	GetRateByDate(ctx context.Context, date string) ([]models.Rate, error)
}

func NewRateService(repo repositories.RateRepositoryInterface) *RateService {
	return &RateService{repo: repo}
}

func (s *RateService) GetAllRates(ctx context.Context) ([]models.Rate, error) {
	return s.repo.GetAllRates(ctx)
}

func (s *RateService) GetRateByDate(ctx context.Context, date string) ([]models.Rate, error) {
	return s.repo.GetRateByDate(ctx, date)
}
