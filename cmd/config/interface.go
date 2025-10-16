package config

import (
	"time"
)

type Argser interface {
	// Server
	GetHost() string

	// Cookie
	GetNameCookie() string
	GetTimeLiveCookie() int

	// Token
	GetTimeLiveToken() time.Duration
	GetSecretToken() []byte

	// Logger
	GetBusinessLog() string
	GetInfraLog() string

	// DB
	GetDB() string

	// Внешний сервис AccrualSystem
	GetSystemAddress() string
	GetMaxConcurrentReq() int
	TimeWaitingResponse() time.Duration
	GetSystemPath() string
	GetTimeTicker() time.Duration
}
