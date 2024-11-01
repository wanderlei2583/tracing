package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

type WeatherService struct {
	apiKey string
}

type weatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{apiKey: apiKey}
}

func (s *WeatherService) GetTemperature(
	ctx context.Context,
	city string,
) (float64, error) {
	ctx, span := otel.Tracer("service-b").Start(ctx, "get_temperature")
	defer span.End()

	url := fmt.Sprintf(
		"http://api.weatherapi.com/v1/current.json?key=%s&q=%s",
		s.apiKey,
		city,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var weather weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return 0, err
	}

	return weather.Current.TempC, nil
}
