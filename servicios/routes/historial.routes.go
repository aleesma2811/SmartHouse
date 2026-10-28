package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *ServicioHandler) GetHistoriaHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	servicio, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Servicio not found"))
		return
	}

	// ?limit=N (por defecto 100, máximo 1000)
	limit := 100
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}

	lecturas, err := h.repo.GetLecturas(params["id"], limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"ServicioID":   servicio.ID,
		"On":           servicio.On,
		"ConsumoTotal": servicio.ConsumoTotal,
		"Lecturas":     lecturas,
	})
}
