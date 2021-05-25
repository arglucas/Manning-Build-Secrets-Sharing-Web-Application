package main

import (
	"github.com/arglucas/secret-app/handlers"
	"log"
	"net/http"
)



// Storage and Global initialization.
// Called on package initialisation (before main etc).
func init() {
	//handlers.InitHandlers()
}

// Start an HTTP Server on :8080 and hook the SecretHandler to it.
func main() {
	addr := ":8080"

	//index := IndexHandler{}
	handler := handlers.SecretHandler{}
	healthHandler := handlers.HealthCheckHandler{}

	//http.HandleFunc("/", Index)
	http.Handle("/healthcheck", handlers.Logger(healthHandler))
	http.Handle("/", handlers.Logger(handler))

	log.Printf("Starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}