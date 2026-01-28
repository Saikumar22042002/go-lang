package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected+"\n" { // json.Encoder adds a newline
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestIsArmstrongFunction(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Positive Case: 153", 153, true},
		{"Positive Case: 371", 371, true},
		{"Positive Case: 9", 9, true},
		{"Zero Case", 0, true},
		{"Negative Case: 100", 100, false},
		{"Negative Case: 123", 123, false},
		{"Negative Number", -153, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isArmstrong(tc.input)
			if result != tc.expected {
				t.Errorf("isArmstrong(%d) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestArmstrongHandler(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		statusCode int
		expected   map[string]interface{}
	}{
		{"Valid Armstrong Number", "/is-armstrong/153", http.StatusOK, map[string]interface{}{"number": 153.0, "isArmstrong": true}},
		{"Valid Non-Armstrong Number", "/is-armstrong/100", http.StatusOK, map[string]interface{}{"number": 100.0, "isArmstrong": false}},
		{"Invalid Number Format", "/is-armstrong/abc", http.StatusBadRequest, nil},
		{"Negative Number", "/is-armstrong/-10", http.StatusBadRequest, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/is-armstrong/{number}", armstrongHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.statusCode)
			}

			if tc.expected != nil {
				var body map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
					t.Fatalf("Could not unmarshal response body: %v", err)
				}
				if body["isArmstrong"] != tc.expected["isArmstrong"] || body["number"] != tc.expected["number"] {
					t.Errorf("handler returned unexpected body: got %v want %v", body, tc.expected)
				}
			}
		})
	}
}
