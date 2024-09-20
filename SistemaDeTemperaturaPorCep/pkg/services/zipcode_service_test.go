package services

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestValidateZipcode(t *testing.T) {
	tests := []struct {
		cep      string
		expected bool
	}{
		{"12345678", true},
		{"1234567", false},
		{"abcdefgh", false},
		{"", false},
	}

	for _, test := range tests {
		result := ValidateZipcode(test.cep)
		assert.Equal(t, test.expected, result)
	}
}

func TestGetLatLonByZipcode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://nominatim.openstreetmap.org/search?postalcode=12345678&country=Brazil&format=json",
		httpmock.NewStringResponder(200, `[{"lat": "-23.5505", "lon": "-46.6333"}]`))

	lat, lon, err := GetLatLonByZipcode("12345678")

	assert.NoError(t, err)
	assert.Equal(t, "-23.5505", lat)
	assert.Equal(t, "-46.6333", lon)
}

func TestGetLatLonByZipcode_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://nominatim.openstreetmap.org/search?postalcode=00000000&country=Brazil&format=json",
		httpmock.NewStringResponder(200, `[]`))

	lat, lon, err := GetLatLonByZipcode("00000000")

	assert.Error(t, err)
	assert.Equal(t, "", lat)
	assert.Equal(t, "", lon)
	assert.EqualError(t, err, "can not find zipcode")
}

func TestGetLatLonByZipcode_HTTPError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://nominatim.openstreetmap.org/search?postalcode=99999999&country=Brazil&format=json",
		httpmock.NewStringResponder(500, "Internal Server Error"))

	lat, lon, err := GetLatLonByZipcode("99999999")

	assert.Error(t, err)
	assert.Equal(t, "", lat)
	assert.Equal(t, "", lon)
	assert.EqualError(t, err, "failed to get zipcode  500")
}
