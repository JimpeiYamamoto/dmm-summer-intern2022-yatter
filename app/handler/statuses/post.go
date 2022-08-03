package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type PostRequest struct {
	Status string
}

type PostReponse struct {
	ID      int64
	Account object.Account
	Status  object.Status
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	account := auth.AccountOf(r)
	status := new(object.Status)
	status.Content = req.Status
	status.AccountID = account.ID
	s := h.app.Dao.Status()
	if err := s.PostStatus(r.Context(), *status); err != nil {
		panic("Invalid object.Status")
	}
	w.Header().Set("Content-Type", "application/json")
	res := PostReponse{ID: status.ID, Account: *account, Status: *status}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	fmt.Println(status)
}
