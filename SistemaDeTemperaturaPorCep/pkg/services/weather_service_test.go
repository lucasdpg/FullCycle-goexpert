package services

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentWeather_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://api.weatherapi.com/v1/current.json?key=testapikey&q=-23.5505,-46.6333&aqi=no",
		httpmock.NewStringResponder(200, `{"current": {"temp_c": 25.0, "temp_f": 77.0}}`))

	lat := "-23.5505"
	lon := "-46.6333"
	apikey := "testapikey"

	temp, err := GetCurrentWeather(lat, lon, apikey)

	assert.NoError(t, err)
	assert.NotNil(t, temp)
	assert.Equal(t, 25.0, temp.Temp_c)
	assert.Equal(t, 77.0, temp.Temp_f)
	assert.Equal(t, 298.15, temp.Temp_k)
}

func TestGetCurrentWeather_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://api.weatherapi.com/v1/current.json?key=testapikey&q=0,0&aqi=no",
		httpmock.NewStringResponder(404, "Not Found"))

	lat := "0"
	lon := "0"
	apikey := "testapikey"

	temp, err := GetCurrentWeather(lat, lon, apikey)

	assert.Error(t, err)
	assert.Nil(t, temp)
	assert.EqualError(t, err, "falha na solicitação: status 404")
}

func TestGetCurrentWeather_BadRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://api.weatherapi.com/v1/current.json?key=wrongapikey&q=-23.5505,-46.6333&aqi=no",
		httpmock.NewStringResponder(400, `{"error": {"message": "Invalid API key."}}`))

	lat := "-23.5505"
	lon := "-46.6333"
	apikey := "wrongapikey"

	temp, err := GetCurrentWeather(lat, lon, apikey)

	assert.Error(t, err)
	assert.Nil(t, temp)
	assert.EqualError(t, err, "falha na solicitação: status 400")
}

func TestGetCurrentWeather_DecodeError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://api.weatherapi.com/v1/current.json?key=testapikey&q=-23.5505,-46.6333&aqi=no",
		httpmock.NewStringResponder(200, `{"current": {"temp_c": "invalid", "temp_f": "invalid"}}`))

	lat := "-23.5505"
	lon := "-46.6333"
	apikey := "testapikey"

	temp, err := GetCurrentWeather(lat, lon, apikey)

	assert.Error(t, err)
	assert.Nil(t, temp)
}
