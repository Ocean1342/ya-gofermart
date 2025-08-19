package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gofermart/internal/auth"
	"gofermart/pkg/common"
	"net/http"
	"strings"
)

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(map[string]interface{}{
			"HANDLER":    "AuthMiddleware",
			"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
		})
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("user not found")))
			return
		}
		token, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			logger.Errorf("error with cut token prefix token: %s %s %s", authHeader, token, found)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("user not found")))
			return
		}
		user, err := h.Auth.GetUserFromToken(auth.Token(token))
		if err != nil || user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Errorf("get user from token error:`%s`", err)
			w.Write([]byte(fmt.Sprintf("user not found")))
			return
		}
		logger.Infof("auth user login:`%s`", user.Login)
		r = r.WithContext(context.WithValue(r.Context(), common.CtxUser, user))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
