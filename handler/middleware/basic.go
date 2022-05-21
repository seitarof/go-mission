package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
)

func BasicAuthentication(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated := BasicAuth(r)
		if !isAuthenticated {
			w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func BasicAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	passwordSha256 := sha256.Sum256([]byte(password))
	p := fmt.Sprintf("%x", passwordSha256[:])
	return username == os.Getenv("BASIC_AUTH_USER_ID") && p == os.Getenv("BASIC_AUTH_PASSWORD")
}
