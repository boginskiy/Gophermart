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
	Config  config.Config
	Auth    auth.Auther
	Logger  logg.Logger
	ResPrep prepar.ResPreper
}

func NewMiddleware(
	config config.Config,
	logger logg.Logger,
	auther auth.Auther,
	resPreper prepar.ResPreper) *Middleware {

	return &Middleware{
		Config:  config,
		Logger:  logger,
		Auth:    auther,
		ResPrep: resPreper}
}

func (m *Middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Этап 1. Регистрация/авторизация пользователя
		hasCredentials := m.Auth.CheckAuthReq(r)
		if hasCredentials {
			next.ServeHTTP(w, r)
			return
		}

		// Этап 2. Аутентификация пользователя
		cookie, err := r.Cookie(m.Config.GetNameCookie())

		// Отсутствуют Cookie. Отправляем сообщение о необходимости регистрации/авторизации
		if err != nil {
			m.ResPrep.ResWithJSON(w, auth.MessNeedRegOrAuth, http.StatusUnauthorized)
			return
		}

		// Присутствуют Cookie. Token просрочен/невалидный
		login, role, id, err := m.Auth.CheckToken(cookie.Value)
		if err != nil {
			m.ResPrep.ResWithJSON(w, auth.MessNeedRegOrAuth, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.CtxUserLogin, login)
		ctx = context.WithValue(ctx, auth.CtxUserRole, role)
		ctx = context.WithValue(ctx, auth.CtxUserID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
