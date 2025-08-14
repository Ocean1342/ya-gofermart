package api

import (
	"github.com/Ocean1342/ya-gofermart/internal/auth"
	"github.com/Ocean1342/ya-gofermart/internal/storage"
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
