package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recover: ", err)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/******************************
 PanicHandler
******************************/
type PanicHandler struct {}

func (PanicHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("--------panic!--------")
}