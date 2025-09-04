package order_processor

import (
	"strconv"
)

type AccrualResponse struct {
	order   string
	status  string
	accrual int
}

func (a *AccrualResponse) GetOrderID() (int, error) {
	return strconv.Atoi(a.order)
}

func (a *AccrualResponse) GetAccrual() int {
	return a.accrual * 100
}

func (a *AccrualResponse) GetStatus() string {
	return a.status
}
