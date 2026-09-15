package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/servicios/db"
	"github.com/servicios/models"
	"github.com/servicios/repository"
	"github.com/servicios/routes"
)

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.Servicio{})

	// Inyección de dependencias
	servicioRepo := repository.NewServicioRepository(db.DB)
	servicioHandler := routes.NewServicioHandler(servicioRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Plugs routes
	r.HandleFunc("/servicios", servicioHandler.GetServiciosHandler).Methods("GET")
	r.HandleFunc("/servicios/{id}", servicioHandler.GetServicioHandler).Methods("GET")
	r.HandleFunc("/servicios", servicioHandler.PostServiciosHandler).Methods("POST")
	r.HandleFunc("/servicios/{id}", servicioHandler.UpdateServicioHandler).Methods("PUT")
	r.HandleFunc("/servicios/{id}", servicioHandler.DeleteServiciosHandler).Methods("DELETE")

	http.ListenAndServe(":4001", r)
}
