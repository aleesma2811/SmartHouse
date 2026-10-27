package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/back/grpcclient"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func validateInmueble(client *grpcclient.InmuebleClient, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(raw))

		var body struct{ InmuebleID uint }
		if err := json.Unmarshal(raw, &body); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		exists, err := client.Exists(r.Context(), body.InmuebleID)
		if err != nil {
			http.Error(w, "no se pudo validar el inmueble por gRPC: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		if !exists {
			http.Error(w, fmt.Sprintf("el inmueble %d no existe", body.InmuebleID), http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
