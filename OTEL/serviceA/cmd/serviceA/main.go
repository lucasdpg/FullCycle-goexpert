package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Cep struct {
	Cep string `json:"cep"`
}

type Temperature struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Temp_K float64 `json:"temp_k"`
}

var tracer trace.Tracer

func main() {
	// Inicializar o OpenTelemetry com Zipkin
	initTracer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(otelMiddleware) // Middleware para traçar todas as requisições

	r.Post("/zipcode", handleFunc)

	http.ListenAndServe(":3030", r)
}

func initTracer() {
	// Configura o Zipkin como o exportador
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("serviceA"),
		)),
	)

	// Configurar o provedor global de traces
	otel.SetTracerProvider(tp)
	tracer = otel.Tracer("zipcode-tracer")
}

func otelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("serviceA").Start(r.Context(), "HTTP "+r.Method)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateZipcode(zipcode string) bool {
	if strings.TrimSpace(zipcode) == "" {
		return false
	}
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

func getServiceB(ctx context.Context, url string) (Temperature, error) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Temperature{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Temperature{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Temperature{}, fmt.Errorf("can not find zipcode: %d", resp.StatusCode)
	}

	var temp Temperature
	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		return Temperature{}, err
	}

	return temp, nil
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("serviceA").Start(r.Context(), "handleFunc")
	defer span.End()

	cep, err := decodeZipcode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	cepok := validateZipcode(cep.Cep)
	if cepok {
		sbUrl := fmt.Sprintf("http://localhost:3000/zipcode-check/%v", cep.Cep)

		// Passa o contexto `ctx` para o `getServiceB`
		temp, err := getServiceB(ctx, sbUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(temp); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}

		fmt.Println(temp)
	} else {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
}
