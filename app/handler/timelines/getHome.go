package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/utils"
)

func (h *handler) GetHome(w http.ResponseWriter, r *http.Request) {
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
	a := h.app.Dao.Account()
	account := auth.AccountOf(r)
	as, err := a.GetFollowings(r.Context(), account.Username, q.Limit)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	s := h.app.Dao.Status()
	statuses := make([]object.Status, 0)
	for _, v := range as {
		fmt.Println(v.Username)
		fmt.Println(v.ID)
		tmp, err := s.GetTimelineHome(r.Context(), q, v.ID)
		statuses = append(statuses, tmp...)
		if err != nil {
			httperror.InternalServerError(w, fmt.Errorf("%w", err))
			return
		}
	}
	ress := make([]Response, 0)
	for _, status := range statuses {
		account, err := a.FindByUserID(r.Context(), status.AccountID)
		if err != nil {
			httperror.InternalServerError(w, fmt.Errorf("%w", err))
		}
		res := Response{
			Id:       status.ID,
			Account:  *account,
			Content:  status.Content,
			CreateAt: status.CreateAt,
		}
		ress = append(ress, res)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ress); err != nil {
		httperror.InternalServerError(w, err)
	}
}
