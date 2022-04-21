package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return username == os.Getenv("BASIC_AUTH_USER_ID") && password == os.Getenv("BASIC_AUTH_PASSWORD")
}