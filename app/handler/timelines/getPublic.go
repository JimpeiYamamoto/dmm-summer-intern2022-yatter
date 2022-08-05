package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/utils"
)

type Response struct {
	Id       int64           `json:"id"`
	Account  object.Account  `json:"account"`
	Content  string          `json:"content"`
	CreateAt object.DateTime `json:"create_at"`
}

func (h *handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	q := object.Query{
		OnlyMedia: r.URL.Query().Get("only_media"),
		MaxID:     r.URL.Query().Get("max_id"),
		SinceID:   r.URL.Query().Get("since_id"),
		Limit:     r.URL.Query().Get("limit"),
	}
	q, err := utils.FixQuery(q)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	s := h.app.Dao.Status()
	statuses, err := s.GetTimelinesPublic(r.Context(), q)
	if err != nil {
		httperror.BadRequest(w, err)
	}
	ress := make([]Response, 0)
	a := h.app.Dao.Account()
	for _, status := range statuses {
		account, err := a.FindByUserID(r.Context(), status.AccountID)
		if err != nil {
			httperror.BadRequest(w, err)
		}
		res := Response{
			Id:       status.ID,
			Account:  *account,
			Content:  status.Content,
			CreateAt: status.CreateAt,
		}
		ress = append(ress, res)
	}
	if err := json.NewEncoder(w).Encode(ress); err != nil {
		httperror.InternalServerError(w, err)
	}
}
