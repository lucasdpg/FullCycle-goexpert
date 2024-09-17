package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucasdpg/FullCycle-goexpert/SistemaDeTemperaturaPorCep/configs"
)

type Temperature struct {
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
	Temp_k float64
}

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
			http.Error(w, "CEP não fornecido", http.StatusBadRequest)
			return
		}

		lat, lon, err := getLatLonByCep(cep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		currentWeather, err := getCurrentWeather(lat, lon, configs.WeatherApiToken)
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

func getLatLonByCep(cep string) (string, string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?postalcode=%s&country=Brazil&format=json", cep)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("falha na solicitação: status %d", resp.StatusCode)
	}

	var jsonData []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return "", "", err
	}

	if len(jsonData) == 0 {
		return "", "", fmt.Errorf("nenhum dado encontrado")
	}

	return jsonData[0].Lat, jsonData[0].Lon, nil
}

func getCurrentWeather(lat, lon, apikey string) (*Temperature, error) {

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s&aqi=no", apikey, lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha na solicitação: status %d", resp.StatusCode)
	}

	var jsonData struct {
		Current struct {
			Temp_c float64 `json:"temp_c"`
			Temp_f float64 `json:"temp_f"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, err
	}

	temp_K := (jsonData.Current.Temp_c + 273.15)

	temp := &Temperature{
		Temp_c: jsonData.Current.Temp_c,
		Temp_f: jsonData.Current.Temp_f,
		Temp_k: temp_K,
	}

	return temp, nil
}
