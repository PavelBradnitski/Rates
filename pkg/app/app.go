package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/PavelBradnitski/Rates/pkg/handlers"
	"github.com/PavelBradnitski/Rates/pkg/models"
	"github.com/PavelBradnitski/Rates/pkg/repositories"
	"github.com/PavelBradnitski/Rates/pkg/services"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
)

const apiURL = "https://api.nbrb.by/exrates/rates?periodicity=0"

func Run() {
	// Подключение к MySql
	dsn := "user_for_migrate:Rn33_io17@tcp(127.0.0.1:3306)/rates_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	rateRepo := repositories.NewRateRepository(db)

	// Создание HTTP сервера
	rateService := services.NewRateService(rateRepo)
	rateHandler := handlers.NewRateHandler(rateService)
	router := gin.Default()
	rateHandler.RegisterRoutes(router)
	go router.Run(":8080")
	now := time.Now()
	sceduleTime := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, time.Local)
	if now.After(sceduleTime) {
		go fetchAndSave()
	}
	c := cron.New()
	_, err = c.AddFunc("0 3 * * *", fetchAndSave)
	if err != nil {
		log.Fatal("Failed scheduled query:", err)
	}
	c.Start()
	// блокировка потока
	select {}
}
func fetchAndSave() {
	dsn := "user_for_migrate:Rn33_io17@tcp(127.0.0.1:3306)/rates_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	Rates, err := fetchRates(apiURL)
	if err != nil {
		log.Fatal("Failed to fetch rates:", err)
	}
	ctx := context.Background()
	rateRepo := repositories.NewRateRepository(db)
	err = rateRepo.AddRates(ctx, Rates)
	if err != nil {
		log.Fatal("Failed to add rates:", err)
	}
}

// Получение курсов из API NBRB
func fetchRates(apiURL string) ([]models.Rate, error) {
	log.Println("Starting Rates fetching...")
	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api returned non-200 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var result []models.Rate
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %w, body: %s", err, string(bodyBytes))
	}
	log.Println("Fetch successful...")
	return result, nil
}
