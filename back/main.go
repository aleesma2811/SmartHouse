package main

import (
	"net/http"

	"github.com/back/db"
	"github.com/back/models"
	"github.com/back/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.Plug{})
	db.DB.AutoMigrate(models.Room{})
	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	// Rooms routes
	r.HandleFunc("/rooms", routes.GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms/{id}", routes.GetRoomHandler).Methods("GET")
	r.HandleFunc("/rooms", routes.PostRoomsHandler).Methods("POST")
	r.HandleFunc("/rooms/{id}", routes.DeleteRoomsHandler).Methods("DELETE")

	// Plugs routes
	r.HandleFunc("/plugs", routes.GetPlugsHandler).Methods("GET")
	r.HandleFunc("/plugs/{id}", routes.GetPlugHandler).Methods("GET")
	r.HandleFunc("/plugs", routes.CreatePlugHandler).Methods("POST")
	r.HandleFunc("/plugs/{id}", routes.UpdatePlugHandler).Methods("PUT")
	r.HandleFunc("/plugs/{id}", routes.DeletePlugsHandler).Methods("DELETE")

	http.ListenAndServe(":4000", r)
}
