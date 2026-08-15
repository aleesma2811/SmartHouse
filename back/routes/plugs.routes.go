package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/back/db"
	"github.com/back/models"
	"github.com/gorilla/mux"
)

func GetPlugsHandler(w http.ResponseWriter, r *http.Request) {
	var plugs []models.Plug
	db.DB.Find(&plugs) // Hacer consulta
	json.NewEncoder(w).Encode(&plugs)
	w.Write([]byte("Get plugs"))

}

func GetPlugHandler(w http.ResponseWriter, r *http.Request) {
	var plug models.Plug
	params := mux.Vars(r)

	db.DB.First(&plug, params["id"])

	if plug.ID == 0 { // Si no existe el room
		w.WriteHeader(http.StatusNotFound) // Error 404
		w.Write([]byte("Room not found"))
		return
	}

	json.NewEncoder(w).Encode(&plug)
}

func CreatePlugHandler(w http.ResponseWriter, r *http.Request) {
	var plug models.Plug

	json.NewDecoder(r.Body).Decode(&plug)

	createdPlug := db.DB.Create(&plug)
	err := createdPlug.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Error 404
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(&plug)
}

func DeletePlugsHandler(w http.ResponseWriter, r *http.Request) {
	var plug models.Plug
	params := mux.Vars(r)
	fmt.Println("ID recibido:", params["id"])

	db.DB.First(&plug, params["id"])
	fmt.Println("Plug encontrado:", plug)

	if plug.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plug not found, invalid ID"))
		return
	}

	db.DB.Unscoped().Delete(&plug)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Plug deleted"))
}

func UpdatePlugHandler(w http.ResponseWriter, r *http.Request) {
	var plug models.Plug
	params := mux.Vars(r)
	fmt.Println("ID a actualizar", params["id"])

	db.DB.First(&plug, params["id"])

	if plug.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plug not found, invalid ID"))
		return
	}

	var updatedData models.Plug
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("PInvalid JSON"))
		return
	}

	db.DB.Model(&plug).Updates(updatedData) // Aplicar cambios
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&plug)

}
