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

	db.DB.AutoMigrate(models.Room{})

	// Inyección de dependencias
	roomRepo := repository.NewRoomRepository(db.DB)
	roomHandler := routes.NewRoomHandler(roomRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Rooms routes
	r.HandleFunc("/rooms", roomHandler.GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms/{id}", roomHandler.GetRoomHandler).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.PostRoomsHandler).Methods("POST")
	r.HandleFunc("/rooms/{id}", roomHandler.DeleteRoomsHandler).Methods("DELETE")

	http.ListenAndServe(":4000", r)
}
