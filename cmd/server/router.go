package server

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/go-chi/chi"
)

type Route struct {
	R                   *chi.Mux
	AuthHandlers        handlers.Hdlrser
	OrdersHandlers      handlers.Hdlrser
	BalanceHandlers     handlers.Hdlrser
	WithdrawalsHandlers handlers.Hdlrser
}

func NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs handlers.Hdlrser) *Route {
	return &Route{
		R:                   chi.NewRouter(),
		AuthHandlers:        authHdlrs,
		OrdersHandlers:      orderHdlrs,
		BalanceHandlers:     balanceHdlrs,
		WithdrawalsHandlers: withdrawalsHdlrs,
	}
}

func (r *Route) Run() http.Handler {
	r.R.Route("/api/user/", func(route chi.Router) {

		r.AuthHandlers.RegisterRoutes(route)        // Авторизация
		r.OrdersHandlers.RegisterRoutes(route)      // Заказы
		r.BalanceHandlers.RegisterRoutes(route)     // Баланс
		r.WithdrawalsHandlers.RegisterRoutes(route) // Выводы средств
	})

	return r.R
}
