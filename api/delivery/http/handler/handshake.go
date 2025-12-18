package handler

import (
	"encoding/json"
	"net/http"
)

func Handshake(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]interface{}{
		"message": "handshake successful!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(resp)
}
