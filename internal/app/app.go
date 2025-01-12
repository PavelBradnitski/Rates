package app

import (
	"fmt"
	"log"
	"os"
	"time"

	scheduler "github.com/PavelBradnitski/Rates/http/client"
	"github.com/PavelBradnitski/Rates/http/server/handler"
	"github.com/PavelBradnitski/Rates/internal/repositories"
	"github.com/PavelBradnitski/Rates/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/robfig/cron/v3"
)

func Run() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	// Подключение к БД
	db, err := repositories.ConnectToDB(dbUser, dbPassword, dbHost, dbPort, dbName)
	defer db.Close()
	// Запуск миграции
	connectionString := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	m, err := migrate.New(
		"file:///Rates/db/migrations",
		connectionString,
	)
	defer m.Close()
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
	rateRepo := repositories.NewRateRepository(db)

	// Создание HTTP сервера
	rateService := services.NewRateService(rateRepo)
	rateHandler := handler.NewRateHandler(rateService)
	router := gin.Default()
	rateHandler.RegisterRoutes(router)
	go router.Run(":8080")
	now := time.Now()
	scheduleTime := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, time.Local)
	// Обработка случая запуска приложения раньше 03:00
	if now.After(scheduleTime) {
		go scheduler.FetchAndSave()
	}
	c := cron.New()
	_, err = c.AddFunc("0 3 * * *", scheduler.FetchAndSave)
	if err != nil {
		log.Fatal("Failed scheduled query:", err)
	}
	c.Start()
	// блокировка потока
	select {}
}
