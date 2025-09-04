package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"net/http"
	"strconv"
	"time"
)

type OrderResponse struct {
	ID         string `json:"number"`
	UserID     int    `json:"-"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "GetOrders",
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
	orders, err := h.Storage.GetOrders(r.Context(), ctxUser.ID)
	if err != nil {
		logger.Errorf("user id:`%d` error get orders: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not define user"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	bytes, err := json.Marshal(hydrate(orders))
	if err != nil {
		logger.Errorf("could not marshal orders")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		logger.Errorf("could not write data to response")
	}
}

func hydrate(orders []*storage.Order) []*OrderResponse {
	res := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		hydratedOrder := &OrderResponse{}
		hydratedOrder.ID = strconv.Itoa(o.ID)
		hydratedOrder.UserID = o.UserID
		hydratedOrder.Status = o.Status
		hydratedOrder.Accrual = o.Accrual
		hydratedOrder.UploadedAt = o.UploadedAt.Format(time.RFC3339)
		hydratedOrder.UpdatedAt = o.UpdatedAt.Time.Format(time.RFC3339)
		res[i] = hydratedOrder
	}
	return res
}
