package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func main() {
	http.HandleFunc("/cotacao", quoteDolar)
	fmt.Println("Server running port 8080")
	http.ListenAndServe(":8080", nil)

}

func dolarSqlite(jsonStr string) error {

	db, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
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
		)
	`)
	if err != nil {
		panic(err)
	}

	var quote Quote
	err = json.Unmarshal([]byte(jsonStr), &quote)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		INSERT INTO USDBRL (
			code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`, quote.USDBRL.Code, quote.USDBRL.Codein, quote.USDBRL.Name, quote.USDBRL.High, quote.USDBRL.Low, quote.USDBRL.VarBid, quote.USDBRL.PctChange, quote.USDBRL.Bid, quote.USDBRL.Ask, quote.USDBRL.Timestamp, quote.USDBRL.CreateDate)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dollar quote recorded in the database")
	return nil
}

func quoteDolar(w http.ResponseWriter, r *http.Request) {

	apiURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Unable to create request", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Request context timed out", http.StatusGatewayTimeout)
			fmt.Println("Timeout context occurred:", err)
		} else {
			http.Error(w, "Unable to connect to dollar quote API", http.StatusInternalServerError)
			fmt.Println("Error occurred:", err)
		}
		return
	}
	defer resp.Body.Close()

	fmt.Println("GET dollar api status recive: ", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

	dolarSqlite(string(body))

}
