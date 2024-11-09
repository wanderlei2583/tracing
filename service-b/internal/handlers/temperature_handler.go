package handlers

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
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

type ErrorResponse struct {
	Message string `json:"message"`
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
		h.sendError(w, "invalid request", http.StatusBadRequest)
		return
	}

	city, err := h.cepService.GetCity(ctx, req.CEP)
	if err != nil {
		switch err {
		case services.ErrInvalidCEP:
			h.sendError(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		case services.ErrCEPNotFound:
			h.sendError(w, "can not find zipcode", http.StatusNotFound)
			return
		default:
			h.sendError(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}
	}

	tempC, err := h.weatherService.GetTemperature(ctx, city)
	if err != nil {
		h.sendError(
			w,
			"error getting temperature",
			http.StatusInternalServerError,
		)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TemperatureHandler) sendError(
	w http.ResponseWriter,
	message string,
	statusCode int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}
