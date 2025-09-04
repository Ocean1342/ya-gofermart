package orderprocessor

import (
	"gofermart/pkg/common"
	"strconv"
)

type AccrualResponse struct {
	Order   string
	Status  string
	Accrual float64
}

func (a *AccrualResponse) GetOrderID() (int, error) {
	return strconv.Atoi(a.Order)
}

func (a *AccrualResponse) GetAccrual() int {
	return common.MoneyFloatToInt(a.Accrual)
}

func (a *AccrualResponse) GetStatus() string {
	return a.Status
}
