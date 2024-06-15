package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	USDBRL struct {
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
	} `json:"USDBRL"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = createTableIfNotExists(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/cotacao", quoteDolarHandler)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS USDBRL (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high REAL,
		low REAL,
		varBid REAL,
		pctChange REAL,
		bid REAL,
		ask REAL,
		timestamp INTEGER,
		create_date TEXT
	)`
	_, err := db.Exec(query)
	return err
}

func quoteDolarHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5000*time.Millisecond)
	defer cancel()

	data, err := fetchQuote(ctx)
	if err != nil {
		http.Error(w, "Unable to fetch quote", http.StatusInternalServerError)
		log.Println("Error fetching quote:", err)
		return
	}

	bid := data.USDBRL.Bid
	if bid == "" {
		http.Error(w, "'bid' field was not found", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{"Dolar": bid}
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)

	dbctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = storeQuote(dbctx, data)
	if err != nil {
		log.Println("Error storing quote:", err)
		panic(err)
	}
}

func fetchQuote(ctx context.Context) (*Quote, error) {
	apiURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("request context timed out: %w", err)
		}
		return nil, fmt.Errorf("unable to connect to dollar quote API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var data Quote
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &data, nil
}

func storeQuote(dbctx context.Context, quote *Quote) error {
	query := `
	INSERT INTO USDBRL (
		code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.ExecContext(dbctx, query,
		quote.USDBRL.Code, quote.USDBRL.Codein, quote.USDBRL.Name,
		quote.USDBRL.High, quote.USDBRL.Low, quote.USDBRL.VarBid,
		quote.USDBRL.PctChange, quote.USDBRL.Bid, quote.USDBRL.Ask,
		quote.USDBRL.Timestamp, quote.USDBRL.CreateDate)
	return err
}
