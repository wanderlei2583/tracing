package handlers

import (
	"encoding/json"
	"net/http"

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

func (h *CEPHandler) HandleCEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "handle_cep_request")
	defer span.End()

	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidCEP(req.CEP) {
		http.Error(w, "invalid CEP", http.StatusUnprocessableEntity)
		return
	}

	resp, err := h.forwardToServiceB(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(resp.Body); err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}
