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
	{Prefix: "/rooms", Target: "http://back:4000"},
	{Prefix: "/servicios", Target: "http://servicios:4001"},
	{Prefix: "/inmuebles", Target: "http://inmuebles:4002"},
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
				log.Printf("LB: %s -> %s", r.URL.Path, route.Prefix)
				proxies[route.Prefix].ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Ruta no encontrada en el load balancer", http.StatusNotFound)
	})

	log.Println("Loadbalancer escuchando en :9000")
	log.Fatal(http.ListenAndServe(":9000", mainHandler))
}
