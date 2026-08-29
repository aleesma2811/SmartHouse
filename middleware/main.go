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
	{Prefix: "/rooms", Target: "http://localhost:4000"},
	{Prefix: "/plugs", Target: "http://localhost:4000"},
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if strings.HasPrefix(r.URL.Path, route.Prefix) {
				proxies[route.Prefix].ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Ruta no encontrada en el gateway", http.StatusNotFound)
	})

	log.Println("Gateway escuchando en :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
