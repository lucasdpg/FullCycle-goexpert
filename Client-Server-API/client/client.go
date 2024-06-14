package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type QuoteResponse struct {
	Dolar string `json:"Dolar"`
}

func getQuote(ctx context.Context) (QuoteResponse, error) {
	apiURL := "http://localhost:8080/cotacao"

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return QuoteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return QuoteResponse{}, fmt.Errorf("request context timed out: %w", err)
		}
		return QuoteResponse{}, fmt.Errorf("unable to connect to dollar quote API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return QuoteResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var quote QuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		return QuoteResponse{}, fmt.Errorf("unable to decode response body: %w", err)
	}

	fmt.Println(quote)
	return quote, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	quote, err := getQuote(ctx)
	if err != nil {
		log.Fatalf("Erro ao obter a cotação: %v", err)
	}

	formattedQuote := fmt.Sprintf("Dolar: {%s}", quote.Dolar)

	err = os.WriteFile("cotacao.txt", []byte(formattedQuote), 0644)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Resposta salva no arquivo cotacao.txt com sucesso!")
}
