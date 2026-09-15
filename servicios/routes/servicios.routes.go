package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/servicios/models"
	"github.com/servicios/repository"
)

type ServicioHandler struct {
	repo repository.ServicioRepository
}

func NewServicioHandler(repo repository.ServicioRepository) *ServicioHandler {
	return &ServicioHandler{repo: repo}
}

func tipoValidado(tipo string) bool {
	switch tipo {
	case models.TipoLuz, models.TipoAgua, models.TipoGas:
		return true
	}
	return false
}

func (h *ServicioHandler) GetServiciosHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId") // Da parámetros del query '/plugs?roomId=10', y se extrae roomId

	var servicios []models.Servicio
	var err error

	if roomID != "" {
		servicios, err = h.repo.GetByRoomID(roomID)
	} else {
		servicios, err = h.repo.GetAll()
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&servicios)
}

func (h *ServicioHandler) GetServicioHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	servicio, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Servicio not found"))
		return
	}
	json.NewEncoder(w).Encode(servicio)
}

func (h *ServicioHandler) PostServiciosHandler(w http.ResponseWriter, r *http.Request) {
	var servicio models.Servicio

	if err := json.NewDecoder(r.Body).Decode(&servicio); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !tipoValidado(servicio.Tipo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tipo inválido: debe ser 'luz', 'agua' o 'gas'"))
		return
	}

	if err := h.repo.Create(&servicio); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&servicio)
}

func (h *ServicioHandler) UpdateServicioHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	servicio, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Servicio not found"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(servicio); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !tipoValidado(servicio.Tipo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tipo inválido: debe ser 'luz', 'agua' o 'gas'"))
		return
	}

	if err := h.repo.Update(servicio); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(servicio)
}

func (h *ServicioHandler) DeleteServiciosHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if _, err := h.repo.GetByID(params["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.repo.Delete(params["id"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sevicio deleted"))
}
