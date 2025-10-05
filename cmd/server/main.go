package server

import (
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/service"
)

// Params need
// SALT_KEY_LEN >> saltLen
// HASH_KEY_LEN >> keyLen

func Start() {
	// // Инициализация main-журнала логирования
	// loggMain := NewLogg()

	// // Инициализация  аргументов
	// args.NewArgs()

	// Services
	userSrv := service.NewUserSrv()

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers()
	balanceHdlrs := handlers.NewBalanceHandlers()
	orderHdlrs := handlers.NewOrdersHandlers()
	authHdlrs := handlers.NewAuthHandlers(userSrv)

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Пуск сервера (args, loggMain)
	NewServer(":8080", "loggMain").Run(router)

}
