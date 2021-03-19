package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

func linkIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if linkID := chi.URLParam(r, "linkID"); linkID != "" {
			ctx := context.WithValue(r.Context(), "linkID", linkID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("unknown page"))
		}
	})
}
