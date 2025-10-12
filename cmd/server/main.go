package server

import (
	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/pkg"
)

func Start(
	args config.Argser,
	appLog logg.Logger,
	infraLog logg.Logger,
	businessLog logg.Logger,
	storeDB store.Dber) {

	// Repository
	repoUsers := repository.NewRepoUsers(args, infraLog, storeDB)
	// repoOrders := repository.NewRepoOrders()
	repo := repository.NewRepo(args, infraLog, repoUsers) !!!!!!!!!

	// Authentification
	JWTServ := auth.NewJWTServ(args, appLog)
	ahCore := auth.NewAhCore(args, appLog)
	auth := auth.NewAuth(ahCore, JWTServ, repo, repoUsers)

	// Checker
	orderChecker := pkg.NewLuna()

	// Services
	// userSrv := service.NewUserSrv(repo, businessLog)
	orderSrv := service.NewOrderSrv(repo, businessLog, orderChecker)

	// Preparation response
	resPrep := prepar.NewResPrep()

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers()
	balanceHdlrs := handlers.NewBalanceHandlers()
	orderHdlrs := handlers.NewOrdersHandlers(orderSrv, resPrep)
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)

	// Middleware
	mdlWare := middleware.NewMiddleware(args, appLog, auth, resPrep)

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Start server
	NewServer(args.GetHost(), appLog).Run(router, mdlWare)

}

// TODO!
// Подключаем БД
// Args доработать
// Repository сделать нормально
// Midlewere доработать
// ВАЖНО! Тестирование
// ВАЖНО! Многопоточность
