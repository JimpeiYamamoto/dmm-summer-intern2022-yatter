package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	query := object.Query{
		OnlyMedia: r.URL.Query().Get("only_media"),
		MaxID:     r.URL.Query().Get("max_id"),
		SinceID:   r.URL.Query().Get("since_id"),
		Limit:     r.URL.Query().Get("limit"),
	}
	s := h.app.Dao.Status()
	statuses, err := s.GetTimelinesPublic(r.Context(), query)
	if err != nil {
		httperror.BadRequest(w, err)
	}
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
	}
}
