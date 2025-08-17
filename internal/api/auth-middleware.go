package api

import (
	"fmt"
	"github.com/Ocean1342/ya-gofermart/internal/auth"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := w.Header().Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("token not found")))
			return
		}
		token, err := h.Auth.GetUserFromToken(auth.Token(authHeader))
		logrus.Infof("fmt: %s %s", token, err)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
