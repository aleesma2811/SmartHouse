package routes

import (
	"encoding/json"
	"net/http"

	"github.com/back/db"
	"github.com/back/models"
	"github.com/gorilla/mux"
)

func GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var rooms []models.Room
	db.DB.Find(&rooms)
	json.NewEncoder(w).Encode(&rooms)
	w.Write([]byte("Get rooms"))
}

func GetRoomHandler(w http.ResponseWriter, r *http.Request) { // Obtener habitación
	var room models.Room
	params := mux.Vars(r)
	// fmt.Println(params["id"])
	db.DB.Preload("Plugs").First(&room, params["id"])

	if room.ID == 0 { // Si no existe el room
		w.WriteHeader(http.StatusNotFound) // Error 404
		w.Write([]byte("Room not found"))
		return
	}

	json.NewEncoder(w).Encode(&room)

}

func PostRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var room models.Room

	json.NewDecoder(r.Body).Decode(&room)

	createdRoom := db.DB.Create(&room)
	err := createdRoom.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(&room)
}

func DeleteRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	params := mux.Vars(r)
	db.DB.First(&room, params["id"])

	if room.ID == 0 { // Si no existe el room
		w.WriteHeader(http.StatusNotFound) // Error 404
		w.Write([]byte("Room not found"))
		return
	}

	db.DB.Unscoped().Delete(&room) // Borrar por completo en la BDD
	w.Write([]byte("Room deleted"))
	w.WriteHeader(http.StatusOK) // 200 OK
}
