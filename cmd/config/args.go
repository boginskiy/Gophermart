package config

import "time"

type Args struct {
}

func NewArgs() *Args {
	return &Args{}
}

func (a *Args) GetHost() string {
	return ":8080"
}

func (a *Args) GetNameCookie() string {
	return "Cookie"
}

func (a *Args) GetTimeLiveCookie() int {
	return 300
}

func (a *Args) GetTimeLiveToken() time.Duration {
	return (20 * time.Second)
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
