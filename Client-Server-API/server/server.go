package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/cotacao", cotacaoDolar)
	http.ListenAndServe(":8080", nil)

}

func cotacaoDolar(w http.ResponseWriter, r *http.Request) {

	apiURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Unable to connect to dollar quote API", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}
