package middleware

import (
	"context"
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type Middleware struct {
	Args config.Argser
	Auth auth.Auther
	Logg logg.Logger
}

func NewMiddleware(argser config.Argser, logger logg.Logger, auther auth.Auther) *Middleware {
	return &Middleware{Args: argser, Logg: logger, Auth: auther}
}

func (m *Middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Этап 1. Пользователь отправил учетные данные в теле запроса
		m.Auth.GetCredentials(r)

		// Этап 2. Аутентификация пользователя
		cookie, err := r.Cookie(m.Args.GetNameCookie())

		// Отсутствуют Cookie. Отправляем сообщение о необходимости регистрации/авторизации
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(auth.MessNeedRegOrAuth)
			return
		}

		// Присутствуют Cookie. Token просрочен/невалидный
		login, role, err := m.Auth.CheckToken(cookie.Value)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(auth.MessNeedRegOrAuth)
			return
		}

		ctx := context.WithValue(r.Context(), auth.CtxLogin, login)
		ctx = context.WithValue(ctx, auth.CtxRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
