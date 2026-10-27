package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/servicios/grpcclient"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func validateRoom(client *grpcclient.RoomClient, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(raw))

		var body struct{ RoomID uint }
		if err := json.Unmarshal(raw, &body); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		exists, err := client.Exists(r.Context(), body.RoomID)
		if err != nil {
			http.Error(w, "no se pudo validar la habitación por gRPC: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		if !exists {
			http.Error(w, fmt.Sprintf("la habitación %d no existe", body.RoomID), http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
