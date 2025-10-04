package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type AuthHandlers struct {
}

func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{}
}

func (ah *AuthHandlers) RegisterRoutes(r chi.Router) {
	r.Post("/register", ah.RegisterUser)
	r.Post("/login", ah.LoginUser)
}

func (ah *AuthHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("RegisterUser"))
}

func (ah *AuthHandlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("LoginUser"))
}
