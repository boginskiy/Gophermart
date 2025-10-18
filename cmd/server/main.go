package server

import (
	"context"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
	"github.com/boginskiy/Gophermart/pkg"
)

func Start(
	args config.Argser,
	appLog logg.Logger,
	infraLog logg.Logger,
	businessLog logg.Logger,
	storeDB store.Dber) {

	// Repository
	repoLoyaltyOrders := repository.NewRepoLoyaltyOrders(args, infraLog, storeDB)
	repoOrders := repository.NewRepoOrders(args, infraLog, storeDB)
	repoUsers := repository.NewRepoUsers(args, infraLog, storeDB)

	// Authentification
	JWTServ := auth.NewJWTServ(args, appLog)
	coreAh := auth.NewCoreAh(args, appLog)
	auth := auth.NewAuth(coreAh, JWTServ, repoUsers)

	// Preparation
	resPrep := prepar.NewResPrep()

	// Chan & Context
	chOrders := make(chan *models.Order, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer close(chOrders)

	// Services
	orderChecker := pkg.NewLuna() // Service проверки номера заказа

	coreSrv := service.NewCoreSrv(args, businessLog, orderChecker)               // coreSrv - сервис с базовым функционалом
	orderSrv := service.NewOrderSrv(chOrders, coreSrv, repoOrders)               // orderSrv - сервис обработки заявок по расчету
	balanceSrv := service.NewBalanceServ(coreSrv, repoOrders, repoLoyaltyOrders) // balanceSrv - сервис обработки бонусов

	service.NewGatewaySrv(ctx, chOrders, coreSrv, repoOrders) // Service прокси для расчета бонусов

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers()
	balanceHdlrs := handlers.NewBalanceHandlers(balanceSrv, resPrep)
	orderHdlrs := handlers.NewOrdersHandlers(orderSrv, resPrep)
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)

	// Middleware
	mdlWare := middleware.NewMiddleware(args, appLog, auth, resPrep)

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Start server
	NewServer(args.GetRunAddress(), appLog).Run(router, mdlWare)

}

// TODO!
// Args доработать
// Err в БД
// Err вообще
// Midlewere доработать
// ВАЖНО! Тестирование
// ВАЖНО! Многопоточность (теория пройти)
