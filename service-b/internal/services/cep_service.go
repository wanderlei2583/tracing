package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Erro bool   `json:"erro,string"`
	CEP  string `json:"cep"`
}

func NewCEPService(baseURL string) *CEPService {
	return &CEPService{baseURL: baseURL}
}

func (s *CEPService) GetCity(ctx context.Context, cep string) (string, error) {
	ctx, span := otel.Tracer("service-b").Start(ctx, "get_city_from_cep")
	defer span.End()

	url := fmt.Sprintf("%s%s/json", s.baseURL, cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", ErrInvalidCEP
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var cepResp viaCEPResponse
	if err := json.Unmarshal(body, &cepResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if cepResp.Erro || cepResp.CEP == "" {
		return "", ErrCEPNotFound
	}

	if cepResp.City == "" {
		return "", ErrCEPNotFound
	}

	return cepResp.City, nil
}
