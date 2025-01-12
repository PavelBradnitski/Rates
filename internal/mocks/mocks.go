package mocks

import (
	"context"

	"github.com/PavelBradnitski/Rates/internal/models"
	"github.com/PavelBradnitski/Rates/internal/repositories"
	"github.com/PavelBradnitski/Rates/internal/services"
	"github.com/stretchr/testify/mock"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetAllRates(ctx context.Context) ([]models.Rate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Rate), args.Error(1)
}

func (m *MockRateService) GetRateByDate(ctx context.Context, date string) ([]models.Rate, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]models.Rate), args.Error(1)
}

var _ services.RateServiceInterface = (*MockRateService)(nil)

type MockRateRepository struct {
	mock.Mock
}

func (m *MockRateRepository) GetAllRates(ctx context.Context) ([]models.Rate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Rate), args.Error(1)
}

func (m *MockRateRepository) GetRateByDate(ctx context.Context, date string) ([]models.Rate, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]models.Rate), args.Error(1)
}

var _ repositories.RateRepositoryInterface = (*MockRateRepository)(nil)
