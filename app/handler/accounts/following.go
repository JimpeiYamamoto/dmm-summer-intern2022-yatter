package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "40"
	}
	a := h.app.Dao.Account()
	as, err := a.GetFollowingUser(r.Context(), username, limit)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(as); err != nil {
		httperror.InternalServerError(w, err)
	}
}
