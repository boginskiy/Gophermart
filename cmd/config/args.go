package config

import "time"

type Args struct {
}

func NewArgs() *Args {
	return &Args{}
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
	return 300
}

// Token
func (a *Args) GetTimeLiveToken() time.Duration {
	return (30 * time.Second)
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
