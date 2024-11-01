package handlers

import (
	"go.opentelemetry.io/otel/trace"

	"service-b/internal/services"
)

type TemperatureHandler struct {
	cepService     *services.CEPService
	weatherService *services.WeatherService
	tracer         trace.Tracer
}

type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
