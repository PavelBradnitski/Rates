package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Rate - структура для хранения данных, получаемых из API
type Rate struct {
	Cur_ID           int     `json:"Cur_ID"`
	Date             string  `json:"Date"`
	Cur_Abbreviation string  `json:"Cur_Abbreviation"`
	Cur_Scale        int     `json:"Cur_Scale"`
	Cur_OfficialRate float64 `json:"Cur_OfficialRate"`
}

func main() {
	// Задайте URL вашего API
	apiURL := "https://api.nbrb.by/exrates/rates?periodicity=0" // Пример API

	fetchAndSaveData(apiURL)
}

// fetchAndSaveData получает данные из API и сохраняет их в файл
func fetchAndSaveData(apiURL string) {
	log.Println("Starting data fetching...")

	// Выполняем запрос к API
	data, err := fetchDataFromAPI(apiURL)
	if err != nil {
		log.Printf("Error fetching data: %v\n", err)
		return
	}

	// Сохраняем данные
	err = saveDataToFile(data, "data.json")
	if err != nil {
		log.Printf("Error saving data: %v\n", err)
		return
	}
	log.Println("Data fetched and saved successfully")
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

// saveDataToFile сохраняет данные в JSON файл
func saveDataToFile(data []Rate, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding json to file: %w", err)
	}

	return nil
}
