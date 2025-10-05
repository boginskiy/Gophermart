package handlers

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/go-chi/chi"
)

type AuthHandlers struct {
	UserSrvcer service.UserSrvcer
}

func NewAuthHandlers(userSrv service.UserSrvcer) *AuthHandlers {
	return &AuthHandlers{UserSrvcer: userSrv}
}

func (ah *AuthHandlers) RegisterRoutes(r chi.Router) {
	r.Post("/register", ah.RegisterUser)
	r.Post("/login", ah.LoginUser)
}

func (ah *AuthHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	dataByte, err := ah.UserSrvcer.Registration(r)

	w.Write([]byte("RegisterUser"))
}

func (ah *AuthHandlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("LoginUser"))
}
