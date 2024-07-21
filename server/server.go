package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeData struct {
	Bid string `json:"bid"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", ExchangeHandler)

	http.ListenAndServe(":8080", mux)
}

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	exchange, err := fetchExchange(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch exchange", http.StatusInternalServerError)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	err = saveExchangeOnDatabase(ctxDB, exchange.Bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// err = saveExchangeOnFile(ctxDB, exchange.Bid)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func fetchExchange(ctx context.Context) (ExchangeData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return ExchangeData{}, fmt.Errorf("Failed to fetch exchange: %w", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ExchangeData{}, fmt.Errorf("Failed to fetch exchange: %w", err)
	}

	defer res.Body.Close()

	var result map[string]map[string]string
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ExchangeData{}, fmt.Errorf("Failed to decode exchange: %w", err)
	}
	exchange := ExchangeData{Bid: result["USDBRL"]["bid"]}

	return exchange, nil
}

func saveExchangeOnDatabase(ctx context.Context, bid string) error {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacao (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"bid" TEXT
	  );`

	_, err = db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return err
	}

	insertSQL := `INSERT INTO cotacao (bid) VALUES (?)`
	_, err = db.ExecContext(ctx, insertSQL, bid)
	return err
}

// func saveExchangeOnFile() {
// 	// Save exchange on file
// }
