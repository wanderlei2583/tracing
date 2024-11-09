package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

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
		http.Error(
			w,
			"error forwarding request",
			http.StatusInternalServerError,
		)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "error reading response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(body); err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func (h *CEPHandler) forwardToServiceB(
	ctx context.Context,
	req CEPRequest,
) (*http.Response, error) {
	_, span := h.tracer.Start(ctx, "forward_to_service_b")
	defer span.End()

	body, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		h.serviceBURL+"/temperature",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}
