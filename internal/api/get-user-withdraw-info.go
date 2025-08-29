package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"net/http"
	"time"
)

type Draw struct {
	Order       int       `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
type DrawsHistory []Draw

func (h *Handler) GetUserWithDrawInfo(w http.ResponseWriter, r *http.Request) {
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
	history, err := h.Storage.GetUserDrawHistory(r.Context(), ctxUser.ID)
	if err != nil {
		logger.Errorf("could not get withdraws. err:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not get withdraws"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}

	if len(history) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	res := make(DrawsHistory, len(history))
	for i, h := range history {
		res[i] = Draw{
			Order:       h.OrderID,
			Sum:         int(h.WithDraw.Int64),
			ProcessedAt: h.TransactionDT,
		}
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		logger.Errorf("could encode withdraws. err:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not get withdraws"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		logger.Errorf("could not write response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
