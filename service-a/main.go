package main

import (
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/telemetry"
)

func main() {
	cleanup := telemetry.InitTracer(
		"service-a",
		"http://zipkin:9411/api/v2/spans",
	)
	defer cleanup()

	r := chi.NewRouter()

	cepHandler := handlers.NewCEPHandler(os.Getenv("SERVICE_B_URL"))
	r.Post(
		"/cep",
		otelhttp.NewHandler(
			http.HandlerFunc(cepHandler.HandleCEP),
			"handle_cep",
		).ServeHTTP,
	)
	log.Fatal(http.ListenAndServe(":8080", r))
}
