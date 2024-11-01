package services

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
