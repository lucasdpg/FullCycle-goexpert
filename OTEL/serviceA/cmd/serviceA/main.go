package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Cep struct {
	Cep string `json:"cep"`
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/zipcode", handleFunc)

	http.ListenAndServe(":3000", r)
}

func validateZipcode(zipcode string) bool {
	regex := regexp.MustCompile(`^\d{8}$`)
	return regex.MatchString(zipcode)
}

func decodeZipcode(body io.ReadCloser) (*Cep, error) {

	var data Cep

	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil

}

func handleFunc(w http.ResponseWriter, r *http.Request) {

	cep, err := decodeZipcode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	cepok := validateZipcode(cep.Cep)
	if cepok {
		fmt.Println(cep.Cep)
	} else {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

}
