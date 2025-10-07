package server

import (
	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/boginskiy/Gophermart/internal/repository"
)

// Params need
// Хеширование пароля
// SALT_KEY_LEN >> saltLen
// HASH_KEY_LEN >> keyLen

// Logger
// названия файлов логирования

// JWT
//  token.SignedString([]byte("SecretKey")) - секретный ключ
//  ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)) - время жизни токена

// Куки
// MaxAge: 300,   // Жива 300 секунд - время жизни
// Name:     name, - имя кук

func Start() {
	// Инициализация main-журнала логирования
	// businessLog := logg.NewLogg("businessLog")
	infraLog := logg.NewLogg("infraLog")
	appLog := logg.NewLogg("appLog")

	// Инициализация  аргументов
	args := config.NewArgs()

	// Repository
	repo := repository.NewRepos(args, infraLog)

	// Authentification
	ahCore := auth.NewAhCore(args, appLog)
	JWTServ := auth.NewJWTServ(args, appLog)
	auth := auth.NewAuth(ahCore, JWTServ, repo)

	// Services
	// userSrv := service.NewUserSrv(repo, businessLog)

	// Handlers
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers()
	balanceHdlrs := handlers.NewBalanceHandlers()
	orderHdlrs := handlers.NewOrdersHandlers()
	authHdlrs := handlers.NewAuthHandlers(auth)

	// Middleware
	mdlWare := middleware.NewMiddleware(appLog, auth)

	// Router
	router := NewRoute(authHdlrs, orderHdlrs, balanceHdlrs, withdrawalsHdlrs)

	// Пуск сервера (args, loggMain)
	NewServer(":8080", appLog).Run(router, mdlWare)

}
