package accrual_system

import "net/http"

//TODO: туть будет реальная имплементация сервиса

type System struct {
}

func New() *System {
	return &System{}
}

func (s *System) OrderProcess(orderID int) *http.Response {
	return &http.Response{}
}
