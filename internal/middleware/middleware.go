package middleware

import (
	"context"
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/prepar"
)

type Middleware struct {
	Args    config.Argser
	Auth    auth.Auther
	Logg    logg.Logger
	ResPrep prepar.ResPreper
}

func NewMiddleware(
	argser config.Argser,
	logger logg.Logger,
	auther auth.Auther,
	resPreper prepar.ResPreper) *Middleware {

	return &Middleware{
		Args:    argser,
		Logg:    logger,
		Auth:    auther,
		ResPrep: resPreper}
}

func (m *Middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Этап 1. Пользователь отправил запрос на регистрацию/авторизацию
		hasCredentials := m.Auth.CheckAuthReq(r)
		if hasCredentials {
			next.ServeHTTP(w, r)
			return
		}

		// Этап 2. Аутентификация пользователя
		cookie, err := r.Cookie(m.Args.GetNameCookie())

		// Отсутствуют Cookie. Отправляем сообщение о необходимости регистрации/авторизации
		if err != nil {
			m.ResPrep.ResWithJson(w, auth.MessNeedRegOrAuth, http.StatusUnauthorized)
			return
		}

		// Присутствуют Cookie. Token просрочен/невалидный
		login, role, err := m.Auth.CheckToken(cookie.Value)
		if err != nil {
			m.ResPrep.ResWithJson(w, auth.MessNeedRegOrAuth, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.CtxLogin, login)
		ctx = context.WithValue(ctx, auth.CtxRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
