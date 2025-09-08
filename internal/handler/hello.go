package handler

import (
	"encoding/json"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Hello Handler",
		"status":  "ok",
	}

	w.Header().Set("Context-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
