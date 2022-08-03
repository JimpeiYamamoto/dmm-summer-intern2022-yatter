package statuses

import (
	"encoding/json"
	"fmt"
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
		panic("ID Must numeric")
	}
	status, err := s.FindById(r.Context(), int64(id))
	if err != nil {
		panic(fmt.Errorf("must existing ID: %w", err))
	}
	fmt.Println("============")
	fmt.Println(status)
	fmt.Println("============")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
	}
}
