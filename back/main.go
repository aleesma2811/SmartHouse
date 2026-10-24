package main

import (
	"context"
	"net/http"
	"time"

	"github.com/back/db"
	"github.com/back/models"
	"github.com/back/repository"
	"github.com/back/routes"
	"github.com/gorilla/mux"
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

	db.DB.AutoMigrate(models.Room{})

	// Inyección de dependencias
	roomRepo := repository.NewRoomRepository(db.DB)
	roomHandler := routes.NewRoomHandler(roomRepo)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/heartbeat", heartbeatHandler).Methods("GET")

	// Rooms routes
	r.HandleFunc("/rooms", roomHandler.GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms/{id}", roomHandler.GetRoomHandler).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.PostRoomsHandler).Methods("POST")
	r.HandleFunc("/rooms/{id}", roomHandler.DeleteRoomsHandler).Methods("DELETE")

	// http.ListenAndServe(":4000", r)
	serve(r, "back", 4000)
}
