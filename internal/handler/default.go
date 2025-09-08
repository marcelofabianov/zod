package handler

import (
	"net/http"
)

func Default(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"message": "OK"}`))
}
