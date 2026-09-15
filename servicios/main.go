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

	db.DB.AutoMigrate(models.Plug{})

	// Inyección de dependencias
	plugRepo := repository.NewPlugRepository(db.DB)
	plugHandler := routes.NewPlugHandler(plugRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Plugs routes
	r.HandleFunc("/plugs", plugHandler.GetPlugsHandler).Methods("GET")
	r.HandleFunc("/plugs/{id}", plugHandler.GetPlugHandler).Methods("GET")
	r.HandleFunc("/plugs", plugHandler.PostPlugsHandler).Methods("POST")
	r.HandleFunc("/plugs/{id}", plugHandler.UpdatePlugHandler).Methods("PUT")
	r.HandleFunc("/plugs/{id}", plugHandler.DeletePlugsHandler).Methods("DELETE")

	http.ListenAndServe(":4001", r)
}
