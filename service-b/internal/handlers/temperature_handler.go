package handlers

import (
	"encoding/json"
	"net/http"
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

func (h *TemperatureHandler) HandleTemperature(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, span := h.tracer.Start(r.Context(), "handle_temperature_request")
	defer span.End()

	var req struct {
		CEP string `json:"cep"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	city, err := h.cepService.GetCity(ctx, req.CEP)
	if err != nil {
		if err == services.ErrInvalidCEP {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if err == services.ErrCEPNotFound {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tempC, err := h.weatherService.GetTemperature(ctx, city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	response := TemperatureResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
