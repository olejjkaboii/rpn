package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/olejjkaboii/rpn/pkg/rpn"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		return
	}

	if req.Expression == "" {
		respondWithError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		return
	}

	result, err := rpn.Calc(req.Expression)
	if err != nil {
		if err.Error() == "деление на ноль" || err.Error() == "неверное число" || err.Error() == "отсутствует закрывающая скобка" {
			respondWithError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, CalculateResponse{Result: result})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, CalculateResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload CalculateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
