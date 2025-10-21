package config

import (
	"time"
)

type Argser interface {
	// Server
	GetRunAddress() string
	// Cookie
	GetNameCookie() string
	GetTimeLiveCookie() int
	// Token
	GetTimeLiveToken() time.Duration
	GetSecretKeyToken() []byte
	// Logger
	GetBusinessLog() string
	GetInfraLog() string
	// DB
	GetDBURI() string
	// BonusCalc
	GetAccrualAddress() string
	GetAccrualMaxReq() int
	GetAccrualWaiteRes() time.Duration
	GetAccrualPath() string
	GetTimeTicker() time.Duration
}
