package main

import (
	"log"
	"net/http"
	"os"
	"service-b/internal/handlers"
	"service-b/internal/services"
	"service-b/internal/telemetry"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	cleanup := telemetry.InitTracer(
		"service-b",
		"http://zipkin:9411/api/v2/spans",
	)
	defer cleanup()

	cepService := services.NewCEPService("https://viacep.com.br/ws/")
	weatherService := services.NewWeatherService(os.Getenv("WEATHER_API_KEY"))

	tempHandler := handlers.NewTemperatureHandler(cepService, weatherService)

	r := chi.NewRouter()
	r.Post(
		"/temperature",
		otelhttp.NewHandler(
			http.HandlerFunc(tempHandler.HandleTemperature),
			"handle_temperature",
		).ServeHTTP,
	)

	log.Fatal(http.ListenAndServe(":8081", r))
}
