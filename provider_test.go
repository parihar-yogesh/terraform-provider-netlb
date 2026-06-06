package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8080")
	if client.address != "http://localhost:8080" {
		t.Errorf("expected address http://localhost:8080, got %s", client.address)
	}
}

func TestExpandStringList(t *testing.T) {
	input := []interface{}{"10.0.0.1:80", "10.0.0.2:80"}
	result := expandStringList(input)
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
	if result[0] != "10.0.0.1:80" {
		t.Errorf("expected 10.0.0.1:80, got %s", result[0])
	}
}

func TestDoRequest_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"name": "test-pool"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.doRequest("GET", "/api/pools/test-pool", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.doRequest("GET", "/api/pools/missing", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// nil response signals resource not found — triggers state removal in Read
	if resp != nil {
		t.Errorf("expected nil response for 404, got %s", string(resp))
	}
}

func TestDoRequest_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json")
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	pool := Pool{Name: "test-pool", LBMethod: "round-robin", Monitor: "http"}
	resp, err := client.doRequest("POST", "/api/pools/test-pool", pool)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestDoRequest_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest("DELETE", "/api/pools/test-pool", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.doRequest("GET", "/api/pools/test", nil)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}
