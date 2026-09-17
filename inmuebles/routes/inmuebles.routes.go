package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inmuebles/models"
	"github.com/inmuebles/repository"
)

type InmuebleHandler struct {
	repo repository.InmuebleRepository
}

func NewInmuebleHandler(repo repository.InmuebleRepository) *InmuebleHandler {
	return &InmuebleHandler{repo: repo}
}

func tipoValidado(tipo string) bool {
	switch tipo {
	case models.TipoCasa, models.TipoDepartamento:
		return true
	}
	return false
}

func (h *InmuebleHandler) GetInmueblesHandler(w http.ResponseWriter, r *http.Request) {
	inmuebles, err := h.repo.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&inmuebles)
}

func (h *InmuebleHandler) GetInmuebleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inmueble, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Inmueble not found"))
		return
	}
	json.NewEncoder(w).Encode(inmueble)
}

func (h *InmuebleHandler) PostInmueblesHandler(w http.ResponseWriter, r *http.Request) {
	var inmueble models.Inmueble

	if err := json.NewDecoder(r.Body).Decode(&inmueble); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !tipoValidado(inmueble.Tipo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tipo inválido: debe ser 'casa' o 'departamento'"))
		return
	}

	if err := h.repo.Create(&inmueble); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&inmueble)
}

func (h *InmuebleHandler) UpdateInmueble(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inmueble, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Servicio not found"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(inmueble); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !tipoValidado(inmueble.Tipo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tipo inválido: debe ser 'casa' o 'departamento'"))
		return
	}

	if err := h.repo.Update(inmueble); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(inmueble)
}

func (h *InmuebleHandler) DeleteInmuebleHandler(w http.ResponseWriter, r *http.Request) {
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
