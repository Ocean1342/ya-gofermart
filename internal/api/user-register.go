package api

import "net/http"

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
	//в транзакции проверять что нет пользователя с таким логином,
	//хешировать пароль и сохранять хешированный пароль
	//создавать(можно параллельно) и сохранять токен

	//на все ошибки отдавать 400
	//проверить уникальность логина, если не уникален, то 409
	//добавить в заголовок полученный токен Authorization: Bearer <token>
}
