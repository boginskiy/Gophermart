package server

import (
	"github.com/boginskiy/Gophermart/internal/handlers"
)

func Start() {
	// // Инициализация main-журнала логирования
	// loggMain := NewLogg()

	// // Инициализация  аргументов
	// args.NewArgs()

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers()
	balanceHdlrs := handlers.NewBalanceHandlers()
	orderHdlrs := handlers.NewOrdersHandlers()
	authHdlrs := handlers.NewAuthHandlers()

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Пуск сервера (args, loggMain)
	NewServer(":8080", "loggMain").Run(router)

}
