package main

import (
	"net/http"

	"github.com/back/db"
	"github.com/back/models"
	"github.com/back/repository"
	"github.com/back/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.Plug{})
	db.DB.AutoMigrate(models.Room{})

	// Inyección de dependencias
	roomRepo := repository.NewRoomRepository(db.DB)
	plugRepo := repository.NewPlugRepository(db.DB)

	roomHandler := routes.NewRoomHandler(roomRepo)
	plugHandler := routes.NewPlugHandler(plugRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Rooms routes
	r.HandleFunc("/rooms", roomHandler.GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms/{id}", roomHandler.GetRoomHandler).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.PostRoomsHandler).Methods("POST")
	r.HandleFunc("/rooms/{id}", roomHandler.DeleteRoomsHandler).Methods("DELETE")

	// Plugs routes
	r.HandleFunc("/plugs", plugHandler.GetPlugsHandler).Methods("GET")
	r.HandleFunc("/plugs/{id}", plugHandler.GetPlugHandler).Methods("GET")
	r.HandleFunc("/plugs", plugHandler.PostPlugsHandler).Methods("POST")
	r.HandleFunc("/plugs/{id}", plugHandler.UpdatePlugHandler).Methods("PUT")
	r.HandleFunc("/plugs/{id}", plugHandler.DeletePlugsHandler).Methods("DELETE")

	http.ListenAndServe(":4000", r)
}
