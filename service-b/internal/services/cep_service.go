package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

var (
	ErrInvalidCEP  = errors.New("invalid CEP")
	ErrCEPNotFound = errors.New("CEP not found")
)

type CEPService struct {
	baseURL string
}

type viaCEPResponse struct {
	City string `json:"localidade"`
	Erro bool   `json:"erro"`
}

func NewCEPService(baseURL string) *CEPService {
	return &CEPService{baseURL: baseURL}
}

func (s *CEPService) GetCity(ctx context.Context, cep string) (string, error) {
	ctx, span := otel.Tracer("service-b").Start(ctx, "get_city_from_cep")
	defer span.End()

	resp, err := http.Get(fmt.Sprintf("%s%s/json", s.baseURL, cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", ErrInvalidCEP
	}

	var cepResp viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepResp); err != nil {
		return "", err
	}

	if cepResp.Erro {
		return "", ErrCEPNotFound
	}

	return cepResp.City, nil
}
