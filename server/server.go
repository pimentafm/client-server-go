package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Exchange struct {
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
}

type ExchangeData struct {
	USDBRL Exchange `json:"USDBRL"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CurrencyExchange)

	http.ListenAndServe(":8080", mux)
}

func CurrencyExchange(w http.ResponseWriter, r *http.Request) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read response: %v\n", err)
	}

	var data ExchangeData
	err = json.Unmarshal(res, &data)

	fmt.Println(data)

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Cotação do Dólar: %s\n", data.USDBRL.Bid))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data.USDBRL.Bid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode JSON: %v\n", err)
	}
}
