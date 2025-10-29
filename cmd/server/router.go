package server

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/go-chi/chi"
)

type Route struct {
	R                   *chi.Mux
	AuthHandlers        handlers.Handler
	OrdersHandlers      handlers.Handler
	BalanceHandlers     handlers.Handler
	WithdrawalsHandlers handlers.Handler
}

func NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs handlers.Handler) *Route {
	return &Route{
		R:                   chi.NewRouter(),
		AuthHandlers:        authHdlrs,
		OrdersHandlers:      orderHdlrs,
		BalanceHandlers:     balanceHdlrs,
		WithdrawalsHandlers: withdrawalsHdlrs,
	}
}

func (r *Route) Run(mv middleware.Mdlwarer) http.Handler {
	// Global Middleware
	r.R.Use(mv.WithAuth)

	r.R.Route("/", func(route chi.Router) {

		// ApiUser
		r.R.Route("/api/user/", func(route chi.Router) {
			r.AuthHandlers.RegisterRoutes(route)        // Авторизация
			r.OrdersHandlers.RegisterRoutes(route)      // Заказы
			r.BalanceHandlers.RegisterRoutes(route)     // Баланс
			r.WithdrawalsHandlers.RegisterRoutes(route) // Выводы средств
		})

	})

	return r.R
}
