package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	account := auth.AccountOf(r)
	a := h.app.Dao.Account()
	targetName := chi.URLParam(r, "username")
	target, err := a.FindByUsername(r.Context(), targetName)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	if err := a.FollowUser(r.Context(), account.ID, target.ID); err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	entity, err := a.GetRelationships(r.Context(), account.ID, target.ID)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
	}
}
