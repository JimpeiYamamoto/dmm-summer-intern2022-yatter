package statuses

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	s := h.app.Dao.Status()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httperror.BadRequest(w, err)
	}
	status, err := s.FindById(r.Context(), int64(id))
	if err != nil {
		httperror.BadRequest(w, err)
	}
	a := h.app.Dao.Account()
	account, err := a.FindByUserID(r.Context(), status.AccountID)
	if err != nil {
		httperror.BadRequest(w, err)
	}
	res := PostReponse{
		Id:       status.ID,
		Account:  *account,
		Content:  status.Content,
		CreateAt: status.CreateAt,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
	}
}
