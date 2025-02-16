package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gogapopp/go-trainee-assignment/internal/http-server/middlewares"
)

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func errorJSONResponse(w http.ResponseWriter, code int, message string) {
	jsonResponse(w, code, map[string]string{"errors": message})
}

func getUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(middlewares.CtxKeyUserID{}).(int)
	if !ok {
		return 0, fmt.Errorf("userID not found in context")
	}
	return userID, nil
}
