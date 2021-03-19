package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/fransoaardi/url-shortener/pkg/api"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	lh := api.NewLinkHandler()

	r.Route("/", func(r chi.Router) {
		r.With(linkIDCtx).Get("/{linkID}", lh.Redirect) // GET /73a14xe
	})

	r.Post("/gen", lh.Generate)

	return r
}
