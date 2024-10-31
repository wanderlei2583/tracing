package handlers

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type CEPHandler struct {
	serviceBURL string
	tracer      trace.Tracer
}

type CEPRequest struct {
	CEP string `json:"cep"`
}

func NewCEPHandler(serviceBURL string) *CEPHandler {
	return &CEPHandler{
		serviceBURL: serviceBURL,
		tracer:      otel.Tracer("service-a"),
	}
}
