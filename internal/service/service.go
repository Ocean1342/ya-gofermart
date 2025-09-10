package service

import (
	"errors"
	"gofermart/internal/storage"
)

type Service struct {
	storage storage.Storage
}

var (
	InternalServerError            = errors.New("internal server")
	PaymentRequiredError           = errors.New("payment required")
	StatusUnprocessableEntityError = errors.New("status unprocessable entity")
)

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
