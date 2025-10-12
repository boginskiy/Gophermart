package config

import "time"

type Argser interface {
	GetNameCookie() string
	GetTimeLiveCookie() int
	GetTimeLiveToken() time.Duration
	GetSecretToken() []byte
	GetHost() string

	// Logger
	GetBusinessLog() string
	GetInfraLog() string

	// DB
	GetDB() string
}
