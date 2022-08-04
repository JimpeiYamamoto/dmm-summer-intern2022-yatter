package timelines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

func (h *handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	var query object.Query
	query.OnlyMedia = r.URL.Query().Get("only_media")
	query.MaxID = r.URL.Query().Get("max_id")
	query.SinceID = r.URL.Query().Get("since_id")
	query.Limit = r.URL.Query().Get("limit")
	fmt.Println("====================================")
	fmt.Println(query.OnlyMedia, query.MaxID, query.SinceID, query.Limit)
	fmt.Println("====================================")
	s := h.app.Dao.Status()
	statuses, err := s.GetTimelinesPublic(r.Context(), query)
	if err != nil {
		panic(fmt.Errorf("select error: %w", err))
	}
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
	}
}
