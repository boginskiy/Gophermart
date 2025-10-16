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

// Удаленный сервис
func (a *Args) GetSystemAddress() string { // ACCRUAL_SYSTEM_ADDRESS
	return "localhost:8081"
}

func (a *Args) GetMaxConcurrentReq() int {
	return runtime.NumCPU()
}

func (a *Args) TimeWaitingResponse() time.Duration {
	return time.Duration(10 * time.Second)
}

func (a *Args) GetTimeTicker() time.Duration {
	return time.Duration(5 * time.Second)
}

func (a *Args) GetSystemPath() string {
	return "/api/orders/"
}

// TODO!
// Params need
// Хеширование пароля
// SALT_KEY_LEN >> saltLen
// HASH_KEY_LEN >> keyLen
