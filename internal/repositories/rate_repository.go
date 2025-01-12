package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/PavelBradnitski/Rates/internal/models"
)

type RateRepository struct {
	db *sql.DB
}
type RateRepositoryInterface interface {
	GetAllRates(ctx context.Context) ([]models.Rate, error)
	GetRateByDate(ctx context.Context, date string) ([]models.Rate, error)
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{db: db}
}

func (r *RateRepository) AddRates(ctx context.Context, rates []models.Rate) error {
	query := `INSERT INTO rates (cur_id, date, cur_abbreviation, cur_scale, cur_official_rate) VALUES (?, ?, ?, ?, ?)`
	for _, rate := range rates {
		parsedDate, _ := time.Parse("2006-01-02T15:04:05", rate.Date)
		_, err := r.db.ExecContext(ctx, query, rate.CurID, parsedDate.Format("2006-01-02"), rate.CurAbbreviation, rate.CurScale, rate.CurOfficialRate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RateRepository) GetAllRates(ctx context.Context) ([]models.Rate, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT ID,cur_id, date, cur_abbreviation, cur_scale, cur_official_rate FROM rates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []models.Rate
	for rows.Next() {
		var rate models.Rate
		err := rows.Scan(&rate.ID, &rate.CurID, &rate.Date, &rate.CurAbbreviation, &rate.CurScale, &rate.CurOfficialRate)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *RateRepository) GetRateByDate(ctx context.Context, date string) ([]models.Rate, error) {
	query := `
		SELECT id,cur_id, date, cur_abbreviation, cur_scale, cur_official_rate 
		FROM rates
		WHERE date = ?
	`
	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	var rates []models.Rate
	for rows.Next() {
		var rate models.Rate
		err := rows.Scan(&rate.ID, &rate.CurID, &rate.Date, &rate.CurAbbreviation, &rate.CurScale, &rate.CurOfficialRate)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	return rates, nil

}
