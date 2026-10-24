package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/servicios/db"
	"github.com/servicios/models"
	"github.com/servicios/repository"
	"github.com/servicios/routes"
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

	db.DB.AutoMigrate(models.Servicio{})

	// Inyección de dependencias
	servicioRepo := repository.NewServicioRepository(db.DB)
	servicioHandler := routes.NewServicioHandler(servicioRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/heartbeat", heartbeatHandler).Methods("GET")

	// Plugs routes
	r.HandleFunc("/servicios", servicioHandler.GetServiciosHandler).Methods("GET")
	r.HandleFunc("/servicios/{id}", servicioHandler.GetServicioHandler).Methods("GET")
	r.HandleFunc("/servicios", servicioHandler.PostServiciosHandler).Methods("POST")
	r.HandleFunc("/servicios/{id}", servicioHandler.UpdateServicioHandler).Methods("PUT")
	r.HandleFunc("/servicios/{id}", servicioHandler.DeleteServiciosHandler).Methods("DELETE")

	serve(r, "servicios", 4001)
}
