package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Temperature struct {
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
	Temp_k float64
}

func GetCurrentWeather(lat, lon, apikey string) (*Temperature, error) {

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
