package config

import (
	"time"
)

type Config interface {
	GetRunAddress() string  // Server
	GetTimeLiveCookie() int // Cookie
	GetNameCookie() string
	GetTimeLiveToken() time.Duration // Token
	GetSecretKeyToken() []byte
	GetBusinessLog() string // Logger
	GetInfraLog() string
	GetDBURI() string                  // DB
	GetAccrualWaiteRes() time.Duration // BonusCalc
	GetTimeTicker() time.Duration
	GetAccrualAddress() string
	GetAccrualPath() string
	GetAccrualMaxReq() int
}
