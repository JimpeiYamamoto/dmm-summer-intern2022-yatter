package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	query := object.Query{
		MaxID:   r.URL.Query().Get("max_id"),
		SinceID: r.URL.Query().Get("since_id"),
		Limit:   r.URL.Query().Get("limit"),
	}
	if query.Limit == "" {
		query.Limit = "40"
	} else {
		limit, err := strconv.ParseUint(query.Limit, 10, 64)
		if err != nil {
			httperror.BadRequest(w, fmt.Errorf("%w", err))
			return
		}
		if limit > 80 {
			query.Limit = "80"
		}
	}
	a := h.app.Dao.Account()
	as, err := a.GetFollowers(r.Context(), username, query)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	if err := json.NewEncoder(w).Encode(as); err != nil {
		httperror.InternalServerError(w, err)
	}
}
