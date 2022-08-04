package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetRelation(w http.ResponseWriter, r *http.Request) {
	account := auth.AccountOf(r)
	a := h.app.Dao.Account()
	targetName := r.URL.Query().Get("username")
	target, err := a.FindByUsername(r.Context(), targetName)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	entity, err := a.GetRelationships(r.Context(), account.ID, target.ID)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	if err := json.NewEncoder(w).Encode(entity); err != nil {
		httperror.InternalServerError(w, err)
	}
}
