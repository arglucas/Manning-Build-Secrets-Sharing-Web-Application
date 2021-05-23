package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// DataStore Data Store for the storing of a mapping from ID to a SecretResponse.
var DataStore map[string]SecretResponse
var mutex *sync.Mutex

type SecretRequest struct{
	PlainText string `json:"plain_text"`
}

type SecretResponse struct{
	Id string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

// Middleware HTTP Handler, that logs the Request.
func logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// Storage and Global initialization.
// Called on package initialisation (before main etc).
func init() {
	// Initialize the data store
	DataStore = make(map[string]SecretResponse)
	mutex = &sync.Mutex{}
	initHandlers()
}

// Start an HTTP Server on :8080 and hook the SecretHandler to it.
func main() {
	addr := ":8080"

	//index := IndexHandler{}
	handler := SecretHandler{}
	healthHandler := HealthCheckHandler{}

	//http.HandleFunc("/", Index)
	http.Handle("/healthcheck", logger(healthHandler))
	http.Handle("/", logger(handler))

	log.Printf("Starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}