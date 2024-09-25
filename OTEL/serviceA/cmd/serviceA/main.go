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

type Temperature struct {
	City   string `json:"city"`
	Temp_C string `json:"temp_C"`
	Temp_F string `json:"temp_F"`
	Temp_K string `json:"temp_K"`
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

func getServiceB(url string) (Temperature, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Temperature{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Temperature{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Temperature{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var temp Temperature
	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		return Temperature{}, err
	}

	return temp, nil

}

func handleFunc(w http.ResponseWriter, r *http.Request) {

	cep, err := decodeZipcode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	cepok := validateZipcode(cep.Cep)
	if cepok {
		sbUrl := fmt.Sprintf("serviceB:3031/find-zipcode/%v", cep.Cep)
		temp, err := getServiceB(sbUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(temp); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}

	} else {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

}
