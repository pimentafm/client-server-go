package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeData struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
	gorm.Model
}

var db *gorm.DB

func main() {
	var err error
	db, err = initDatabase()
	if err != nil {
		fmt.Println(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", ExchangeHandler)

	http.ListenAndServe(":8080", mux)
}

func initDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&ExchangeData{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	exchange, err := fetchExchange(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	err = saveExchangeOnDatabase(ctxDB, exchange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchange)
}

func fetchExchange(ctx context.Context) (ExchangeData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return ExchangeData{}, fmt.Errorf("failed to fetch exchange: %w", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ExchangeData{}, fmt.Errorf("failed to fetch exchange: %w", err)
	}

	defer res.Body.Close()

	var result map[string]map[string]string
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ExchangeData{}, fmt.Errorf("failed to decode exchange: %w", err)
	}

	exchange := ExchangeData{
		Code:       result["USDBRL"]["code"],
		Codein:     result["USDBRL"]["codein"],
		Name:       result["USDBRL"]["name"],
		High:       result["USDBRL"]["high"],
		Low:        result["USDBRL"]["low"],
		VarBid:     result["USDBRL"]["varBid"],
		PctChange:  result["USDBRL"]["pctChange"],
		Bid:        result["USDBRL"]["bid"],
		Ask:        result["USDBRL"]["ask"],
		Timestamp:  result["USDBRL"]["timestamp"],
		CreateDate: result["USDBRL"]["create_date"],
	}

	return exchange, nil
}

func saveExchangeOnDatabase(ctx context.Context, exchange ExchangeData) error {
	db.WithContext(ctx).AutoMigrate(&ExchangeData{})

	result := db.WithContext(ctx).Create(&exchange)
	if result.Error != nil {
		return fmt.Errorf("failed to save exchange on database: %w", result.Error)
	}

	return nil
}
