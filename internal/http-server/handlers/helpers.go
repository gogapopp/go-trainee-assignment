package handlers

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func errorJSONResponse(w http.ResponseWriter, code int, message string) {
	jsonResponse(w, code, map[string]string{"errors": message})
}
