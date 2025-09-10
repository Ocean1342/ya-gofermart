package service

import (
	"gofermart/internal/storage"
	"net/http"
)

type Serviceable interface {
	UserWithDraw(r *http.Request, user *storage.User, withdrawRequest WithdrawRequest) error
}
