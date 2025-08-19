package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gofermart/pkg/common"
	"io"
	"net/http"
)

type UserAuthRequest struct {
	Login    string `json:"login,required"`
	Password string `json:"password,required"`
}

func (h *Handler) UserAuth(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "UserAuth",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	var userAuthRequest UserAuthRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to read request body. err: %v", err)))
		return
	}
	err = json.Unmarshal(bodyBytes, &userAuthRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to read request body. err: %v", err)))
		return
	}
	user, err := h.Storage.GetUserByLogin(r.Context(), userAuthRequest.Login)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("user not found. err: %v", err)))
		return
	}
	if !h.Auth.CompareHashAndPassword(userAuthRequest.Password, user.Password) {
		logger.Errorf("401: %s %s", userAuthRequest, user)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := h.Auth.CreateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not create token. err: %v", err)))
		return
	}
	logger.Debugf("user login:`%s` password:`%s` token:`%s` created", user.Login, user.Password, token)
	w.Header().Set(common.AuthorizationHeaderName, fmt.Sprintf("Bearer %s", string(token)))

}
