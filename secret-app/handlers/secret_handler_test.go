package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecretFetch(t *testing.T) {
	expectedOutput := `{"data":"My super secret123"}`
	DataStore["c616584ac64a93aafe1c16b6620f5bcd"] = SecretResponse{Id: "c616584ac64a93aafe1c16b6620f5bcd", Data: "My super secret123"}


	req, err := http.NewRequest("GET", "/c616584ac64a93aafe1c16b6620f5bcd", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	secretHandler := SecretHandler{}
	handler := http.Handler(secretHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) != expectedOutput {
		t.Errorf("handler returned unexpected body: got %v want %v",
			strings.TrimSpace(rr.Body.String()), expectedOutput)
	}

	// Check Second Fetch returns Empty
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	if strings.TrimSpace(rr.Body.String()) != EmptySecretResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			strings.TrimSpace(rr.Body.String()), EmptySecretResponse)
	}

	// Check get to '/' with no ID gives an error
	// Check Second Fetch returns Empty
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	if strings.TrimSpace(rr.Body.String()) != EmptySecretResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			strings.TrimSpace(rr.Body.String()), EmptySecretResponse)
	}
}

func TestSecretStore(t *testing.T) {
	inputBody := []byte(`{"plain_text":"My super secret123"}`)
	expectedOutput := `{"id":"c616584ac64a93aafe1c16b6620f5bcd"}`
	//DataStore["c616584ac64a93aafe1c16b6620f5bcd"] = SecretResponse{Id: "c616584ac64a93aafe1c16b6620f5bcd", Data: "My super secret123"}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(inputBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	secretHandler := SecretHandler{}
	handler := http.Handler(secretHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) != expectedOutput {
		t.Errorf("handler returned unexpected body: got %v want %v",
			strings.TrimSpace(rr.Body.String()), expectedOutput)
	}

	// POST to '/' with no body
	req, err = http.NewRequest("POST", "/", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if strings.TrimSpace(rr.Body.String()) != EmptyJsonResponse {
		t.Errorf("handler returned unexpected body: got %v want %v",
			strings.TrimSpace(rr.Body.String()), EmptyJsonResponse)
	}
}

func TestSecretPut(t *testing.T) {
	req, err := http.NewRequest("PUT", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	secretHandler := SecretHandler{}
	handler := http.Handler(secretHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}