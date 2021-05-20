package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

// Data Store for the storing of a mapping from ID to a SecretResponse.
var DataStore map[string]SecretResponse

type SecretRequest struct{
	PlainText string `json:"plain_text"`
}

type SecretResponse struct{
	Id string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

var EmptySecretResponse string

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
		json.NewEncoder(w).Encode(sr)

		// Add Data and store it
		sr.Data = d.PlainText
		DataStore[sr.Id] = sr
		log.Printf("Stored Id: %s", sr.Id)

	case "GET": // Retrieve Secret
		id := html.EscapeString(r.URL.Path)[1:] // remote slash at beginning
		sr, exists := DataStore[id]
		if !exists { // if it doesn't exist return empty data
			log.Printf("No value for Id: %s", id)
			fmt.Fprintf(w, EmptySecretResponse)
			return
		}
		log.Printf("Retrieved Id: %s", sr.Id)
		sr.Id = "" // remove the Id from the response
		json.NewEncoder(w).Encode(sr)
		// Remove the key now it has been returned to the caller
		delete(DataStore, id)
	}
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
	// Prepare a Empty SecretResponse for easier to read responses.
	EmptySecretResponse = "{\"data\": \"\"}"
}

// Generate an MD5 hash of a supplied string, return as a hex string
func md5hex(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// Start an HTTP Server on :8080 and hook the SecretHandler to it.
func main() {
	addr := ":8080"

	//index := IndexHandler{}
	handler := SecretHandler{}

	//http.HandleFunc("/", Index)
	http.Handle("/", logger(handler))

	log.Printf("Starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
