package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucasdpg/FullCycle-goexpert/SistemaDeTemperaturaPorCep/configs"
	"github.com/lucasdpg/FullCycle-goexpert/SistemaDeTemperaturaPorCep/pkg/services"
)

func main() {

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/temperature-by-cep/{cep}", func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")

		if cep == "" {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if !services.ValidateZipcode(cep) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		lat, lon, err := services.GetLatLonByZipcode(cep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		currentWeather, err := services.GetCurrentWeather(lat, lon, configs.WeatherApiToken)
		if err != nil {
			http.Error(w, "Error Weather: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Define the Content-Type header as application/json
		w.Header().Set("Content-Type", "application/json")

		// Marshal the temperature data to JSON and write it to the response
		json.NewEncoder(w).Encode(currentWeather)
	})

	http.ListenAndServe(":3000", r)
}
