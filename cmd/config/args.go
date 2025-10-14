package config

import (
	"runtime"
	"time"

	"github.com/boginskiy/Gophermart/internal/logg"
)

type Args struct {
	Logg logg.Logger
}

func NewArgs(logger logg.Logger) *Args {
	return &Args{Logg: logger}
}

// Host
func (a *Args) GetHost() string {
	return ":8080"
}

// Cookie
func (a *Args) GetNameCookie() string {
	return "auth_cookie"
}

func (a *Args) GetTimeLiveCookie() int {
	return 3000
}

// Token
func (a *Args) GetTimeLiveToken() time.Duration {
	return (3000 * time.Second)
}

func (a *Args) GetSecretToken() []byte {
	return []byte("SecretKey")
}

// Logger
func (a *Args) GetInfraLog() string {
	return "infraLog"
}

func (a *Args) GetBusinessLog() string {
	return "businessLog"
}

// DataBase
func (a *Args) GetDB() string {
	return "postgres://username:userpassword@localhost:5432/gophermartdb?sslmode=disable"
}

// ACCRUAL_SYSTEM_ADDRESS
func (a *Args) GetACCRUAL_SYSTEM_ADDRESS() string {
	return "localhost:8081"
}

func (a *Args) GetMaxConcurrentReq() int {
	return runtime.GOMAXPROCS(0)
}

// TODO!
// Params need
// Хеширование пароля
// SALT_KEY_LEN >> saltLen
// HASH_KEY_LEN >> keyLen
