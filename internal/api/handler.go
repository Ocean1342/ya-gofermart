package api

import (
	"gofermart/internal/auth"
	order_processor "gofermart/internal/order-processor"
	"gofermart/internal/storage"
)

type Handler struct {
	Storage        storage.Storage
	Auth           auth.Auth
	OrderProcessor *order_processor.OrderProcessor
}

func New(storage storage.Storage, auth auth.Auth, processor *order_processor.OrderProcessor) *Handler {
	return &Handler{
		Storage:        storage,
		Auth:           auth,
		OrderProcessor: processor,
	}
}
