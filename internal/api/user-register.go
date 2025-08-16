package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

/*
TODO:

	структура запроса
	структура ответа
	мидл вар для проверки токенов на определённые ручки
*/
type UserRegisterRequest struct {
	Login    string `json:"login,required"`
	Password string `json:"password,required"`
}

func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "UserRegister",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	var userRegisterRequest UserRegisterRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to read request body. err: %v", err)))
		return
	}
	err = json.Unmarshal(bodyBytes, &userRegisterRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to read request body. err: %v", err)))
		return
	}

	//в транзакции проверять что нет пользователя с таким логином,
	hashPassword, err := h.Auth.CreateHash(userRegisterRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("could not create hash for password")))
		return
	}
	user, err := h.Storage.CreateUser(r.Context(), userRegisterRequest.Login, hashPassword)
	//TODO: проверить ошибку sql на дубль, вдруг что то происходит
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		logger.Errorf("user login:`%s` already exists. err:`%s`", userRegisterRequest.Login, err)
		w.Write([]byte(fmt.Sprintf("user login:`%s` already exists", userRegisterRequest.Login)))
		return
	}
	logger.Debugf("user login:`%s` password:`%s` created", user.Login, user.Password)
	//сохранять токен

	//на все ошибки отдавать 400
	//проверить уникальность логина, если не уникален, то 409
	//добавить в заголовок полученный токен Authorization: Bearer <token>
}
