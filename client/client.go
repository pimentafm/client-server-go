package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Exchange struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 900*time.Second)
	defer cancel()

	exchange, err := fetchExchange(ctx)
	if err != nil {
		log.Fatalf("Erro ao obter cotação: %v\n", err)
	}

	err = saveExchangeToFile(exchange.Bid)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação: %v\n", err)
	}
}

func fetchExchange(ctx context.Context) (Exchange, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return Exchange{}, fmt.Errorf("Erro ao criar requisição: %v\n", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Exchange{}, fmt.Errorf("Erro ao fazer requisição: %v\n", err)
	}
	defer res.Body.Close()

	var exchange Exchange
	err = json.NewDecoder(res.Body).Decode(&exchange)
	if err != nil {
		return Exchange{}, fmt.Errorf("Erro ao decodificar resposta: %v\n", err)
	}

	return exchange, nil
}

func saveExchangeToFile(bid string) error {
	content := fmt.Sprintf("Dólar: %s", bid)

	f, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("Erro ao criar arquivo: %v\n", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("Erro ao escrever no arquivo: %v\n", err)
	}

	return nil

}
