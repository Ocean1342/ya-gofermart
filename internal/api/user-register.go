package api

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	//смаршалить в структуру, если ошибка, то 400
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

	tx, err := h.Storage.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not start transaction. err: %v", err)))
		return
	}
	fmt.Println(tx)
	//в транзакции проверять что нет пользователя с таким логином,

	//хешировать пароль и сохранять хешированный пароль
	//создавать(можно параллельно) и сохранять токен

	//на все ошибки отдавать 400
	//проверить уникальность логина, если не уникален, то 409
	//добавить в заголовок полученный токен Authorization: Bearer <token>
}
