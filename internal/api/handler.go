package api

import (
	"gofermart/internal/auth"
	order_processor "gofermart/internal/order-processor"
	"gofermart/internal/service"
	"gofermart/internal/storage"
)

type Handler struct {
	Storage        storage.Storage
	Auth           auth.Auth
	OrderProcessor *order_processor.OrderProcessor
	Service        service.Serviceable
}

func New(storage storage.Storage, auth auth.Auth, processor *order_processor.OrderProcessor, service service.Serviceable) *Handler {
	return &Handler{
		Storage:        storage,
		Auth:           auth,
		OrderProcessor: processor,
		Service:        service,
	}
}
