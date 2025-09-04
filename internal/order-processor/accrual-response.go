package orderprocessor

import (
	"strconv"
)

type AccrualResponse struct {
	Order   string
	Status  string
	Accrual int
}

func (a *AccrualResponse) GetOrderID() (int, error) {
	return strconv.Atoi(a.Order)
}

func (a *AccrualResponse) GetAccrual() int {
	return a.Accrual * 100
}

func (a *AccrualResponse) GetStatus() string {
	return a.Status
}
