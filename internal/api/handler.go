package api

import (
	"gofermart/internal/auth"
	"gofermart/internal/storage"
)

type Handler struct {
	Storage storage.Storage
	Auth    auth.Auth
}

func New(storage storage.Storage, auth auth.Auth) *Handler {
	return &Handler{
		Storage: storage,
		Auth:    auth,
	}
}
