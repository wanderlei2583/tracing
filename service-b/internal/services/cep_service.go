package services

import "errors"

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
