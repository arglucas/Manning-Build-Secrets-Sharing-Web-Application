package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

// Handler format to explain how to add a chained handler and how it looks conceptually.
//
//type IndexHandler struct{}
// Old handler format
// func Index(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
// }
// Chainable Handler Format
//func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//}

var HealthCheckResponse string
var EmptySecretResponse string

// General Initialization of handlers for default responses, errors etc.
func initHandlers() {
	// Prepare a health check response response.
	HealthCheckResponse = "{\"Status\": \"OK\"}"

	// Prepare a Empty SecretResponse for easier to read responses.
	EmptySecretResponse = "{\"data\": \"\"}"

}

// HealthCheckHandler provides a health checking route to confirm the server is working

type HealthCheckHandler struct{}
func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, HealthCheckResponse)
}

// SecretHandler provides a way for storing a secret and retrieving it exactly once.
// POST - Body -- {"plain_text": "secret"} -- Returns {"id": "hexstring"}
// GET  - URI  -- /hexstring -- Returns {"data": "secret"}
// If the URI passed to get isn't in the data store -- Returns {"data": ""}

type SecretHandler struct{}
func (h SecretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "test handler\n")
	switch r.Method {
	case "POST": // Store Secret
		var d SecretRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&d); err != nil {
			log.Println("Error: Request in wrong format. " + err.Error())
			return
		}
		//fmt.Println(d.PlainText, md5hex(d.PlainText), reflect.TypeOf(d))
		sr := SecretResponse{Id: md5hex(d.PlainText)}
		// Send to stream before adding the data to the structure
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sr)

		// Add Data and store it
		mutex.Lock()
		sr.Data = d.PlainText
		DataStore[sr.Id] = sr
		mutex.Unlock()
		log.Printf("Stored Id: %s", sr.Id)

	case "GET": // Retrieve Secret
		id := html.EscapeString(r.URL.Path)[1:] // remote slash at beginning
		mutex.Lock()
		sr, exists := DataStore[id]
		mutex.Unlock()
		if !exists { // if it doesn't exist return empty data
			log.Printf("No value for Id: %s", id)
			fmt.Fprintf(w, EmptySecretResponse)
			return
		}
		log.Printf("Retrieved Id: %s", sr.Id)
		sr.Id = "" // remove the Id from the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sr)
		// Remove the key now it has been returned to the caller
		mutex.Lock()
		delete(DataStore, id)
		mutex.Unlock()
	}
}
