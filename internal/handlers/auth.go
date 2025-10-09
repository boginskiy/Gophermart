package handlers

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/go-chi/chi"
)

type AuthHandlers struct {
	Auther auth.Auther
}

func NewAuthHandlers(auther auth.Auther) *AuthHandlers {
	return &AuthHandlers{Auther: auther}
}

func (ah *AuthHandlers) RegisterRoutes(r chi.Router) {
	r.Post("/register", ah.RegisterUser)
	r.Post("/login", ah.LoginUser)
}

func (ah *AuthHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	dataByte, cookie, err := ah.Auther.Registration(r)

	// Логин, пароль введены некорректно
	if err == auth.ErrLoginPasswordIsBad {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Введенный логин занят другим пользователем
	if err == auth.ErrLoginNotUnic {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}

	// Все ошибки создания нового пользователя
	if err == auth.ErrCreateUser {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataByte)
}

func (ah *AuthHandlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	dataByte, cookie, err := ah.Auther.Authentication(r)

	// Логин, пароль введены некорректно
	if err == auth.ErrLoginPasswordIsBad {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Логин не найден среди зарегистрированных пользователей || Невалидный пароль
	if err == auth.ErrLogindNotFound || err == auth.ErrPasswordNotValid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	// Все ошибки создания нового пользователя
	if err == auth.ErrCreateUser {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataByte)
}
