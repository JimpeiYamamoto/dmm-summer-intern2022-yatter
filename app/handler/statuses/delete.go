package statuses

import (
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s := h.app.Dao.Status()
	if err := s.DeleteStatus(r.Context(), id); err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
	}
	w.Header().Set("Content-Type", "application/json")
}
