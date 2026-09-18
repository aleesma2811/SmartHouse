package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
)

type Route struct {
	Prefix  string   // inicio del path
	Targets []string // A donde va la petición
	counter uint64
}

// Tabla de ruteo
var routes = []Route{
	{Prefix: "/rooms", Targets: []string{"http://back-1:4000", "http://back-2:4000"}},
	{Prefix: "/servicios", Targets: []string{"http://servicios-1:4001", "http://servicios-2:4001"}},
	{Prefix: "/inmuebles", Targets: []string{"http://inmuebles-1:4002", "http://inmuebles-2:4002"}},
}

func (route *Route) nextTarget() string {
	n := atomic.AddUint64(&route.counter, 1)
	index := int(n % uint64(len(route.Targets)))
	return route.Targets[index]
}

func main() {
	proxies := make(map[string]*httputil.ReverseProxy)

	for _, route := range routes {
		for _, target := range route.Targets {
			parsedURL, err := url.Parse(target)
			if err != nil {
				log.Fatal(err)
			}
			proxies[target] = httputil.NewSingleHostReverseProxy(parsedURL)
		}
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := range routes {
			route := &routes[i]
			if strings.HasPrefix(r.URL.Path, route.Prefix) {
				target := route.nextTarget()
				log.Printf("LB: %s -> %s", r.URL.Path, target)
				proxies[target].ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Ruta no encontrada en el load balancer", http.StatusNotFound)
	})

	log.Println("Loadbalancer escuchando en :9000")
	log.Fatal(http.ListenAndServe(":9000", mainHandler))
}
