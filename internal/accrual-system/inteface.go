package accrualSystem

import "net/http"

type AccrualServiceResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual,omitempty"`
}

type AccrualService interface {
	OrderProcess(orderID int) *http.Response
}
