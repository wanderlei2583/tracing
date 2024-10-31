package main

func main() {
	cleanup := telemetry.InitTelemetry(
		"service-a",
		"http://zipkin:9411/api/v2/spans",
	)
	defer cleanup()

	r := chi.NewRouter()
}
