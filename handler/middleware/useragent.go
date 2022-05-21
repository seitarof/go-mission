package middleware

import (
	"context"
	"net/http"

	ua "github.com/mileusna/useragent"
)

func UserAgent(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get a UserAgent and OS name
		ua := ua.Parse(r.UserAgent())
		r = r.WithContext(context.WithValue(r.Context(), "OS", ua.OS))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
