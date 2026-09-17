package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inmuebles/db"
	"github.com/inmuebles/models"
	"github.com/inmuebles/repository"
	"github.com/inmuebles/routes"
)

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.Inmueble{})

	// Inyección de dependencias
	inmuebleRepo := repository.NewInmuebleRepository(db.DB)
	inmuebleHandler := routes.NewInmuebleHandler(inmuebleRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Plugs routes
	r.HandleFunc("/inmuebles", inmuebleHandler.GetInmueblesHandler).Methods("GET")
	r.HandleFunc("/inmuebles/{id}", inmuebleHandler.GetInmuebleHandler).Methods("GET")
	r.HandleFunc("/inmuebles", inmuebleHandler.PostInmueblesHandler).Methods("POST")
	r.HandleFunc("/inmuebles/{id}", inmuebleHandler.UpdateServicioHandler).Methods("PUT")
	r.HandleFunc("/inmuebles/{id}", inmuebleHandler.DeleteInmuebleHandler).Methods("DELETE")

	log.Println("Servicio 'inmuebles' escuchando en :4002")
	http.ListenAndServe(":4002", r)
}
