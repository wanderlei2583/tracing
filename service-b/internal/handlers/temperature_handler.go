package handlers

import (
	"service-b/internal/services"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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

func NewTemperatureHandler(
	cep *services.CEPService,
	weather *services.WeatherService,
) *TemperatureHandler {
	return &TemperatureHandler{
		cepService:     cep,
		weatherService: weather,
		tracer:         otel.Tracer("service-b"),
	}
}
