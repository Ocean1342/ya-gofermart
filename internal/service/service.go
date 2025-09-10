package service

import (
	"errors"
	"gofermart/internal/storage"
)

type Service struct {
	storage storage.Storage
}

var (
	ErrInternalServer            = errors.New("internal server")
	ErrPaymentRequired           = errors.New("payment required")
	ErrStatusUnprocessableEntity = errors.New("status unprocessable entity")
)

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
