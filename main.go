package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Rate struct {
	CurID           int     `json:"Cur_ID"`
	Date            string  `json:"Date"`
	CurAbbreviation string  `json:"Cur_Abbreviation"`
	CurScale        int     `json:"Cur_Scale"`
	CurOfficialRate float64 `json:"Cur_OfficialRate"`
}

var db *sql.DB

func main() {
	dsn := "user_for_migrate:Rn33_io17@tcp(127.0.0.1:3306)/rates_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	apiURL := "https://api.nbrb.by/exrates/rates?periodicity=0"

	Rates, err := fetchRates(apiURL)
	if err != nil {
		log.Fatal("Failed to fetch rates:", err)
	}
	// TODO: к addRatesToDB нужно добавить ресивер db
	//err = addRatesToDB(Rates)
	query := `INSERT INTO rates (cur_id, date, cur_abbreviation, cur_scale, cur_official_rate) VALUES (?, ?, ?, ?, ?)`
	for _, rate := range Rates {
		parsedDate, _ := time.Parse("2006-01-02T15:04:05", rate.Date)
		_, err := db.Exec(query, rate.CurID, parsedDate.Format("2006-01-02"), rate.CurAbbreviation, rate.CurScale, rate.CurOfficialRate)
		if err != nil {
			log.Fatal("Failed to insert rate:", err)
		}
	}
	if err != nil {
		log.Fatal("Failed to insert rates into database:", err)
	}

	http.HandleFunc("/rates", getAllRatesHandler)
	http.HandleFunc("/rate", getRateByDateHandler)

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// fetchAndSaveData получает данные из API и сохраняет их в файл
func fetchRates(apiURL string) ([]Rate, error) {
	log.Println("Starting Rates fetching...")

	// Выполняем запрос к API
	data, err := fetchDataFromAPI(apiURL)
	if err != nil {
		log.Printf("Error fetching data: %v\n", err)
		return nil, err
	}
	log.Println("Rates fetched successfully")
	return data, nil
}

// fetchDataFromAPI получает данные из API
func fetchDataFromAPI(apiURL string) ([]Rate, error) {
	client := http.Client{Timeout: 10 * time.Second} // Установка таймаута

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

	var result []Rate
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %w, body: %s", err, string(bodyBytes))
	}

	return result, nil
}

// func addRatesToDB(rates []Rate) error {

// }

func getAllRatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT cur_id, date, cur_abbreviation, cur_scale, cur_official_rate FROM rates")
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var rates []Rate
	for rows.Next() {
		var rate Rate
		if err := rows.Scan(&rate.CurID, &rate.Date, &rate.CurAbbreviation, &rate.CurScale, &rate.CurOfficialRate); err != nil {
			http.Error(w, "Failed to scan data", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		rates = append(rates, rate)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}

func getRateByDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Date parameter is required", http.StatusBadRequest)
		return
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	query := `SELECT cur_id, date, cur_abbreviation, cur_scale, cur_official_rate FROM rates WHERE date = ?`
	row := db.QueryRow(query, date)

	var rate Rate
	if err := row.Scan(&rate.CurID, &rate.Date, &rate.CurAbbreviation, &rate.CurScale, &rate.CurOfficialRate); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No rate found for the given date", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		log.Println("QueryRow error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}
