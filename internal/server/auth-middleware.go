package server

import (
	"net/http"
)

func authable(next http.Handler) http.Handler {
	//TODO: validate token
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
