package config

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/caarlos0/env"
)

type ArgsENV struct {
	Logg logg.Logger

	RunAddress      string `env:"RUN_ADDRESS"`            //
	NameCookie      string `env:"NAME_COOKIE"`            //
	TimeLiveCookie  int    `env:"TIME_LIVE_COOKIE"`       //
	SecretKeyToken  string `env:"SECRET_KEY_TOKEN"`       //
	TimeLiveToken   int    `env:"TIME_LIVE_TOKEN"`        //
	InfraLog        string `env:"INFRA_LOG_FILE"`         //
	BusinessLog     string `env:"BUSINESS_LOG_FILE"`      //
	DBURI           string `env:"DATABASE_URI"`           //
	TimeTicker      int    `env:"TIME_TICKER"`            //
	AccrualAddress  string `env:"ACCRUAL_SYSTEM_ADDRESS"` //
	AccrualMaxReq   int    `env:"ACCRUAL_MAX_REQUEST"`    //
	AccrualWaiteRes int    `env:"ACCRUAL_WAITE_RESPONSE"` //
	AccrualPath     string `env:"ACCRUAL_PATH"`           //
}

func NewArgsENV(logger logg.Logger) *ArgsENV {
	args := &ArgsENV{Logg: logger}
	args.ParseFlags()
	return args
}

func (e *ArgsENV) ParseFlags() {
	err := env.Parse(e)
	if err != nil {
		e.Logg.RaiseError("ArgsENV>ParseFlags>Parse", err)
	}

	// Default
	valueStr := strings.TrimSpace(os.Getenv("RUN_ADDRESS"))
	if len(valueStr) == 0 {
		e.RunAddress = ":8080"
	}

	valueStr = strings.TrimSpace(os.Getenv("NAME_COOKIE"))
	if len(valueStr) == 0 {
		e.NameCookie = "auth_cookie"
	}

	valueStr = strings.TrimSpace(os.Getenv("TIME_LIVE_COOKIE"))
	if len(valueStr) == 0 {
		e.TimeLiveCookie = 3000
	}

	valueStr = strings.TrimSpace(os.Getenv("SECRET_KEY_TOKEN"))
	if len(valueStr) == 0 {
		e.SecretKeyToken = "Ld5pS4Gw"
	}

	valueStr = strings.TrimSpace(os.Getenv("TIME_LIVE_TOKEN"))
	if len(valueStr) == 0 {
		e.TimeLiveToken = 3000
	}

	valueStr = strings.TrimSpace(os.Getenv("INFRA_LOG_FILE"))
	if len(valueStr) == 0 {
		e.InfraLog = "infraLog"
	}

	valueStr = strings.TrimSpace(os.Getenv("BUSINESS_LOG_FILE"))
	if len(valueStr) == 0 {
		e.BusinessLog = "businessLog"
	}

	valueStr = strings.TrimSpace(os.Getenv("DATABASE_URI"))
	if len(valueStr) == 0 {
		e.DBURI = "postgres://username:userpassword@localhost:5432/gophermartdb?sslmode=disable"
	}

	valueStr = strings.TrimSpace(os.Getenv("TIME_TICKER"))
	if len(valueStr) == 0 {
		e.TimeTicker = 1
	}

	valueStr = strings.TrimSpace(os.Getenv("ACCRUAL_SYSTEM_ADDRESS"))
	if len(valueStr) == 0 {
		e.AccrualAddress = "localhost:8081"
	}

	valueStr = strings.TrimSpace(os.Getenv("ACCRUAL_MAX_REQUEST"))
	if len(valueStr) == 0 {
		e.AccrualMaxReq = runtime.NumCPU()
	}

	valueStr = strings.TrimSpace(os.Getenv("ACCRUAL_WAITE_RESPONSE"))
	if len(valueStr) == 0 {
		e.AccrualWaiteRes = 1
	}

	valueStr = strings.TrimSpace(os.Getenv("ACCRUAL_PATH"))
	if len(valueStr) == 0 {
		e.AccrualPath = "/api/orders/"
	}
}

func (e *ArgsENV) GetRunAddress() string {
	return e.RunAddress
}

func (e *ArgsENV) GetNameCookie() string {
	return e.NameCookie
}

func (e *ArgsENV) GetTimeLiveCookie() int {
	return e.TimeLiveCookie
}

func (e *ArgsENV) GetTimeLiveToken() time.Duration {
	return time.Duration(time.Duration(e.TimeLiveToken) * time.Second)
}

func (e *ArgsENV) GetSecretKeyToken() []byte {
	return []byte(e.SecretKeyToken)
}

func (e *ArgsENV) GetInfraLog() string {
	return e.InfraLog
}

func (e *ArgsENV) GetBusinessLog() string {
	return e.BusinessLog
}

func (e *ArgsENV) GetDBURI() string {
	return e.DBURI
}

func (e *ArgsENV) GetAccrualAddress() string {
	return e.AccrualAddress
}

func (e *ArgsENV) GetAccrualMaxReq() int {
	return e.AccrualMaxReq
}

func (e *ArgsENV) GetAccrualWaiteRes() time.Duration {
	return time.Duration(time.Duration(e.AccrualWaiteRes) * time.Second)
}

func (e *ArgsENV) GetTimeTicker() time.Duration {
	return time.Duration(time.Duration(e.TimeTicker) * time.Second)
}

func (e *ArgsENV) GetAccrualPath() string {
	return e.AccrualPath
}
