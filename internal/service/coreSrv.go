package service

import (
	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type CoreSrv struct {
	Args       config.Argser
	Logg       logg.Logger
	OrderCheck OrderChecker // OrderCheck is verification by the Luna algorithm
}

func NewCoreSrv(a config.Argser, l logg.Logger, orderCheck OrderChecker) *CoreSrv {
	return &CoreSrv{
		Args:       a,
		Logg:       l,
		OrderCheck: orderCheck,
	}
}
