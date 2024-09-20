package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func GetLatLonByZipcode(cep string) (string, string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?postalcode=%s&country=Brazil&format=json", cep)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get zipcode  %d", resp.StatusCode)
	}

	var jsonData []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return "", "", err
	}

	if len(jsonData) == 0 {
		return "", "", fmt.Errorf("can not find zipcode")
	}

	return jsonData[0].Lat, jsonData[0].Lon, nil
}

func ValidateZipcode(cep string) bool {
	regex := regexp.MustCompile(`^\d{8}$`)
	return regex.MatchString(cep)
}
