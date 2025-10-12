package handlers

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/go-chi/chi"
)

type AuthHandlers struct {
	Auther  auth.Auther
	ResPrep prepar.ResPreper
}

func NewAuthHandlers(auther auth.Auther, resPreper prepar.ResPreper) *AuthHandlers {
	return &AuthHandlers{Auther: auther, ResPrep: resPreper}
}

func (ah *AuthHandlers) RegisterRoutes(r chi.Router) {
	r.Post("/register", ah.RegisterUser)
	r.Post("/login", ah.LoginUser)
}

func (ah *AuthHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	dataByte, cookie, err := ah.Auther.Registration(r)

	// Логин, пароль введены некорректно
	if err == auth.ErrLoginPasswordIsBad {
		ah.ResPrep.ResWithJson(w, []byte(err.Error()), http.StatusBadRequest)
		return
	}

	// Введенный логин занят другим пользователем
	if err == auth.ErrLoginNotUnic {
		ah.ResPrep.ResWithJson(w, []byte(err.Error()), http.StatusConflict)
		return
	}

	// Все ошибки создания нового пользователя
	if err == auth.ErrCreateUser {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ah.ResPrep.ResWithJsonAndCookie(w, dataByte, cookie, http.StatusOK)
}

func (ah *AuthHandlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	dataByte, cookie, err := ah.Auther.Authentication(r)

	// Логин, пароль введены некорректно
	if err == auth.ErrLoginPasswordIsBad || err == auth.ErrLoginPasswordIsBad2 {
		ah.ResPrep.ResWithJson(w, []byte(err.Error()), http.StatusBadRequest)
		return
	}

	// Логин не найден среди зарегистрированных пользователей || Невалидный пароль
	if err == auth.ErrLogindNotFound || err == auth.ErrPasswordNotValid {
		ah.ResPrep.ResWithJson(w, []byte(err.Error()), http.StatusUnauthorized)
		return
	}

	// Все ошибки создания нового пользователя
	if err == auth.ErrCreateUser {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ah.ResPrep.ResWithJsonAndCookie(w, dataByte, cookie, http.StatusOK)
}
