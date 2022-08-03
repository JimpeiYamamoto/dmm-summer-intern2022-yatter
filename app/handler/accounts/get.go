package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	a := h.app.Dao.Account()
	account, err := a.FindByUsername(r.Context(), username)
	if err != nil {
		panic("Must existing Username")
	}
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
	}
}
