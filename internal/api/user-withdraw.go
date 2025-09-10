package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"io"
	"net/http"
)

func (h *Handler) UserWithdraw(w http.ResponseWriter, r *http.Request) {
	var withdrawRequest service.WithdrawRequest
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "UserWithdraw",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	ctxUser, ok := (r.Context().Value(common.CtxUser)).(*storage.User)
	if !ok || ctxUser == nil {
		logger.Errorf("could not extract user from context")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("could not define user"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("user id:`%d` could not read request body: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("could not read request body"))
		if err != nil {
			logger.Errorf("could not read request body")
		}
		return
	}
	err = json.Unmarshal(bytes, &withdrawRequest)
	if err != nil {
		logger.Errorf("user id:`%d` could not unmarshal body: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(" could not unmarshal body request"))
		if err != nil {
			logger.Errorf(" could not unmarshal body request")
		}
		return
	}
	err = h.Service.UserWithDraw(r, ctxUser, withdrawRequest)
	if err != nil {
		if errors.Is(err, service.ErrInternalServer) {
			logger.Errorf("service error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(fmt.Sprintf("could not withdraw: %s", err.Error())))
			if err != nil {
				logger.Errorf("could not write data to response")
			}
			return
		} else if errors.Is(err, service.ErrPaymentRequired) {
			logger.Errorf("service error: %s", err)
			w.WriteHeader(http.StatusPaymentRequired)
			_, err = w.Write([]byte(fmt.Sprintf("could not withdraw: %s", err.Error())))
			if err != nil {
				logger.Errorf("could not write data to response")
			}
			return
		} else if errors.Is(err, service.ErrStatusUnprocessableEntity) {
			logger.Errorf("service error: %s", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err = w.Write([]byte(fmt.Sprintf("could not withdraw: %s", err.Error())))
			if err != nil {
				logger.Errorf("could not write data to response")
			}
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("service error: %s", err)
			_, err = w.Write([]byte("unknown error"))
			if err != nil {
				logger.Errorf("could not write data to response")
			}
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
