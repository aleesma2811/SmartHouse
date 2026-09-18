package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Route struct {
	Prefix string // inicio del path
	Target string // A donde va la petición
}

// Tabla de ruteo
var routes = []Route{
	{Prefix: "/rooms", Target: "http://loadbalancer:9000"},
	{Prefix: "/servicios", Target: "http://loadbalancer:9000"},
	{Prefix: "/inmuebles", Target: "http://loadbalancer:9000"},
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	proxies := make(map[string]*httputil.ReverseProxy)

	for _, route := range routes {
		target, err := url.Parse(route.Target) // Convierte la ruta a un objeto *url.URL
		if err != nil {
			log.Fatal(err)
		}
		proxies[route.Prefix] = httputil.NewSingleHostReverseProxy(target) // Crea un objeto proxy
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if strings.HasPrefix(r.URL.Path, route.Prefix) {
				proxies[route.Prefix].ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Ruta no encontrada en el gateway", http.StatusNotFound)
	})

	log.Println("Gateway escuchando en :8000")
	log.Fatal(http.ListenAndServe(":8000", corsMiddleware(mainHandler)))
}
