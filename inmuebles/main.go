package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/inmuebles/db"
	"github.com/inmuebles/models"
	"github.com/inmuebles/repository"
	"github.com/inmuebles/routes"
)

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := db.DB.DB()
	if err == nil {
		// Si se obtuvo la conexión, hacer un ping a Posgres con límite de 2s
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		err = sqlDB.PingContext(ctx)
	}
	if err != nil {
		http.Error(w, "DB unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Back alive"))
}

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.Inmueble{})

	// Inyección de dependencias
	inmuebleRepo := repository.NewInmuebleRepository(db.DB)
	inmuebleHandler := routes.NewInmuebleHandler(inmuebleRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/heartbeat", heartbeatHandler).Methods("GET")

	// Plugs routes
	r.HandleFunc("/inmuebles", inmuebleHandler.GetInmueblesHandler).Methods("GET")
	r.HandleFunc("/inmuebles/{id}", inmuebleHandler.GetInmuebleHandler).Methods("GET")
	r.HandleFunc("/inmuebles", inmuebleHandler.PostInmueblesHandler).Methods("POST")
	r.HandleFunc("/inmuebles/{id}", inmuebleHandler.UpdateInmueble).Methods("PUT")
	r.HandleFunc("/inmuebles/{id}", inmuebleHandler.DeleteInmuebleHandler).Methods("DELETE")

	log.Println("Servicio 'inmuebles' escuchando en :4002")
	serve(r, "inmuebles", 4002)
}
