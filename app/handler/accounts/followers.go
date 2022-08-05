package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/utils"

	"github.com/go-chi/chi"
)

func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	q := object.Query{
		MaxID:   r.URL.Query().Get("max_id"),
		SinceID: r.URL.Query().Get("since_id"),
		Limit:   r.URL.Query().Get("limit"),
	}
	q, err := utils.FixQuery(q)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	a := h.app.Dao.Account()
	as, err := a.GetFollowers(r.Context(), username, q)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(as); err != nil {
		httperror.InternalServerError(w, err)
	}
}
