package api

import "net/http"

func (h *Handler) UserAuth(w http.ResponseWriter, r *http.Request) {
	//проверить заголовок полученный токен Authorization: Bearer <token>
}
