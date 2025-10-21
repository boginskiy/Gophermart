package server

import (
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/middleware"
)

type Serv struct {
	args   config.Argser
	Logger logg.Logger
}

func NewServer(argser config.Argser, logger logg.Logger) *Serv {
	return &Serv{
		args:   argser,
		Logger: logger,
	}
}

func (s *Serv) Run(router Router, mv middleware.Mdlwarer) {
	s.Logger.RaiseFatal(
		"server is not running",
		http.ListenAndServe(s.args.GetRunAddress(), router.Run(mv)))
}
