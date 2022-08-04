package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Post("/", h.Create)
	r.Get("/{username}", h.Get)
	r.With(auth.Middleware(app)).Post("/{username}/follow", h.Follow)
	r.With(auth.Middleware(app)).Get("/relationships", h.GetRelation)
	r.Get("/{username}/following", h.Following)
	return r
}
