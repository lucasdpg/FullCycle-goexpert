package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucasdpg/FullCycle-goexpert/OTEL/serviceB/configs"
	otelpkg "github.com/lucasdpg/FullCycle-goexpert/OTEL/serviceB/pkg/otel"
	"github.com/lucasdpg/FullCycle-goexpert/OTEL/serviceB/pkg/services"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func main() {

	// Initialize OpenTelemetry tracer
	cleanup := otelpkg.InitTracer()
	defer cleanup()

	WeatherApiToken := os.Getenv("WEATHER_API_TOKEN")
	fmt.Println(WeatherApiToken)
	if WeatherApiToken == "" {
		fmt.Println("Getting WEATHER_API_TOKEN in .env file")
		configs, err := configs.LoadConfig(".")
		if err != nil {
			panic(err)
		}
		WeatherApiToken = configs.WeatherApiToken
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "request")
	})

	r.Get("/zipcode-check/{cep}", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("serviceB").Start(r.Context(), "handle-zipcode-check")
		defer span.End()

		cep := chi.URLParam(r, "cep")

		if cep == "" {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if !services.ValidateZipcode(cep) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		lat, lon, err := services.GetLatLonByZipcode(ctx, cep) // Passar ctx
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		currentWeather, err := services.GetCurrentWeather(ctx, lat, lon, WeatherApiToken) // Passar ctx
		if err != nil {
			http.Error(w, "Error Weather: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentWeather)
	})

	http.ListenAndServe(":3000", otelhttp.NewHandler(r, "server"))
}
