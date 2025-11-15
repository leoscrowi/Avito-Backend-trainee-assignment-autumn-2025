package utils

import (
	"encoding/json"
	"net/http"
)

func WriteHeader(w http.ResponseWriter, statusCode int, item interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(item)
	if err != nil {
		return
	}
}
