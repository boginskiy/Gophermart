package server

import (
	"context"

	conf "github.com/boginskiy/Gophermart/cmd/config"
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
	config conf.Config,
	appLog logg.Logger,
	infraLog logg.Logger,
	businessLog logg.Logger,
	storeDB store.DataBase) {

	// Repository
	repoLoyaltyOrders := repository.NewRepoLoyaltyOrders(config, infraLog, storeDB)
	repoOrders := repository.NewRepoOrders(config, infraLog, storeDB)
	repoUsers := repository.NewRepoUsers(config, infraLog, storeDB)

	// Authentification
	JWTServ := auth.NewJWTServ(config, appLog)
	auth := auth.NewAuth(config, appLog, repoUsers, JWTServ)

	// Preparation
	resPrep := prepar.NewResPrep()

	// Chan & Context
	chOrders := make(chan *models.Order, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer close(chOrders)

	// Services
	orderChecker := pkg.NewLuna()                                                                          // orderChecker - сервис  проверки номера заказа
	orderSrv := service.NewOrderSrv(chOrders, config, businessLog, repoOrders, orderChecker)               // orderSrv - сервис обработки заявок по расчету
	balanceSrv := service.NewBalanceServ(config, businessLog, repoOrders, repoLoyaltyOrders, orderChecker) // balanceSrv - сервис обработки бонусов
	withdrawalsSrv := service.NewWithdrawalsSrv(config, businessLog, repoLoyaltyOrders)                    // withdrawalsSrv - сервис обработки использованных бонусов
	service.NewGatewaySrv(ctx, chOrders, config, businessLog, repoOrders)                                  // gatewaySrv - сервис прокси для расчета бонусов

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers(withdrawalsSrv, resPrep)
	balanceHdlrs := handlers.NewBalanceHandlers(balanceSrv, resPrep)
	orderHdlrs := handlers.NewOrdersHandlers(orderSrv, resPrep)
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)

	// Middleware
	mdlWare := middleware.NewMiddleware(config, appLog, auth, resPrep)

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Start server
	NewServer(config, appLog).Run(router, mdlWare)

}

// TODO!
// Args доработать
// Err в БД
// Midlewere доработать. Сжатие
// Документацию!
// Тестирование!

// По теории
// Error, Многопоточность ...
