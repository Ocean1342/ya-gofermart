package orderprocessor

import (
	"gofermart/pkg/common"
	"math"
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
	round := math.Round(a.Accrual*100) / 100
	return common.MoneyFloatToInt(round)
}

func (a *AccrualResponse) GetStatus() string {
	return a.Status
}
