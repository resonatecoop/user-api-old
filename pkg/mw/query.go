package mw

import (
  "net/http"
  "context"
  "github.com/go-pg/pg/urlvalues"
)

func WithURLQuery(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()
		q := urlvalues.Values{"limit": query["limit"], "page": query["page"]}
		ctx = context.WithValue(ctx, "query", q)
		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}
