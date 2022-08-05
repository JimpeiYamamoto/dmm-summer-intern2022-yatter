package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetHome(w http.ResponseWriter, r *http.Request) {
	query := object.Query{
		OnlyMedia: r.URL.Query().Get("only_media"),
		MaxID:     r.URL.Query().Get("max_id"),
		SinceID:   r.URL.Query().Get("since_id"),
		Limit:     r.URL.Query().Get("limit"),
	}
	l, err := strconv.ParseUint(query.Limit, 10, 64)
	if err != nil {
		query.Limit = "40"
	} else if l > 80 {
		query.Limit = "80"
	}
	a := h.app.Dao.Account()
	account := auth.AccountOf(r)
	as, err := a.GetFollowingUser(r.Context(), account.Username, query.Limit)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf("%w", err))
		return
	}
	s := h.app.Dao.Status()
	statuses := make([]object.Status, 0)
	for _, v := range as {
		fmt.Println(v.Username)
		fmt.Println(v.ID)
		tmp, err := s.GetTimelineHome(r.Context(), query, v.ID)
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
