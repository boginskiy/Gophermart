package server

import (
	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/repository"
)

// Params need
// Хеширование пароля
// SALT_KEY_LEN >> saltLen
// HASH_KEY_LEN >> keyLen

// JWT
//  token.SignedString([]byte("SecretKey")) - секретный ключ
//  ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)) - время жизни токена

// Куки
// MaxAge: 300,   // Жива 300 секунд - время жизни
// Name:     name, - имя кук

func Start() {

	// Args & Logger
	appLog := logg.NewLogg("appLog")
	args := config.NewArgs()

	businessLog := logg.NewLogg(args.GetBusinessLog())
	infraLog := logg.NewLogg(args.GetInfraLog())

	// Repository
	repo := repository.NewRepos(args, infraLog)

	// Authentification
	JWTServ := auth.NewJWTServ(args, appLog)
	ahCore := auth.NewAhCore(args, appLog)
	auth := auth.NewAuth(ahCore, JWTServ, repo)

	// Services
	// userSrv := service.NewUserSrv(repo, businessLog)

	// Preparation Response
	resPrep := prepar.NewResPrep()

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers()
	balanceHdlrs := handlers.NewBalanceHandlers()
	orderHdlrs := handlers.NewOrdersHandlers()
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)

	// Middleware
	mdlWare := middleware.NewMiddleware(args, appLog, auth, resPrep)

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Start server
	NewServer(args.GetHost(), appLog).Run(router, mdlWare)

	// Clouse somethings
	defer businessLog.Clouse()
	defer infraLog.Clouse()
	defer appLog.Clouse()
}
