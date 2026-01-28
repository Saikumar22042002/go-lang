package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// isArmstrong checks if a number is an Armstrong number.
// An Armstrong number is a number that is the sum of its own digits each raised to the power of the number of digits.
func isArmstrong(num int) bool {
	if num < 0 {
		return false
	}

	numStr := strconv.Itoa(num)
	numDigits := len(numStr)
	var sum float64
	originalNum := num

	for originalNum > 0 {
		digit := originalNum % 10
		sum += math.Pow(float64(digit), float64(numDigits))
		originalNum /= 10
	}

	return int(sum) == num
}

// healthHandler responds with a simple health check message.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// armstrongHandler handles requests to check for Armstrong numbers.
func armstrongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	numStr, ok := vars["number"]
	if !ok {
		http.Error(w, "Number not provided", http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid number format", http.StatusBadRequest)
		return
	}

	if num < 0 {
		http.Error(w, "Number must be a non-negative integer", http.StatusBadRequest)
		return
	}

	result := isArmstrong(num)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"number":     num,
		"isArmstrong": result,
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/is-armstrong/{number}", armstrongHandler).Methods("GET")

	port := "8080"
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
