package statuses

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type PostRequest struct {
	Status string
}

type PostReponse struct {
	Id       int64           `json:"id:"`
	Account  object.Account  `json:"account"`
	Content  string          `json:"content"`
	CreateAt object.DateTime `json:"create_at"`
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	account := auth.AccountOf(r)
	status := object.Status{
		Content:   req.Status,
		AccountID: account.ID,
	}
	s := h.app.Dao.Status()
	if err := s.PostStatus(r.Context(), &status); err != nil {
		httperror.BadRequest(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	res := PostReponse{
		Id:       status.ID,
		Account:  *account,
		Content:  status.Content,
		CreateAt: status.CreateAt}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
