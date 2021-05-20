package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMD5(t *testing.T) {
	//buf := new(bytes.Buffer)
	//
	//displayGreetings(buf)
	//
	//expectedOutput := "Hello!\nWorld!\n"
	//gotOutput := buf.String()
	//if gotOutput != expectedOutput {
	//	t.Fatalf("Expected: %s, Got: %s\n", expectedOutput, gotOutput )
	//}
	gotOutput := md5hex("My super secret123")

	expectedOutput := "c616584ac64a93aafe1c16b6620f5bcd"
	if gotOutput != expectedOutput {
		t.Fatalf("Expected: %s, Got: %s\n", expectedOutput, gotOutput )
	}
}

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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
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
}