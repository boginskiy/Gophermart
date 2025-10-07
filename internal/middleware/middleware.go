package middleware

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type Middleware struct {
	Auther auth.Auther
	Logger logg.Logger
}

func NewMiddleware(logger logg.Logger, auther auth.Auther) *Middleware {
	return &Middleware{Logger: logger, Auther: auther}
}

func (m *Middleware) WithAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Регистрация пользователя
		// dataByte, err := m.Auther.Registration(r)

		// if err == auth.ErrLoginPasswordIsBad {
		// 	// ВозвратBadRequest
		// 	return
		// }

		// Registration(r)
		// ReAuthentication(r)
		// Authentication(r)

		next.ServeHTTP(w, r)
	})
}

// Логика работы:
// Новый клиент должен зарегистрироваться. Попадает в БД с данными
// Выдаем ему токен, куки, для прохождения дальнейшей аутентификации

// Далее клиент зарегистрирован, но у него закончился просрочился токен
// Проверяем что пользак наш, выдаем ему свежий токен

// Пользак пришел без токена/без регистрации/битые данные, минимальный функционал
// Просьба ему зарегаться и получить полный функционал
//

// Доп, после трех попыток неверного ввода пароля, просьба о восстановлении пароля

// Авто токен действует в течении суток, далее его продляем только через повторный ввод пароля

// TODO!
// func(http.Handler) http.Handler
