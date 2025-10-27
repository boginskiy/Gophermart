package server

import (
	"net/http"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/middleware"
)

type Serv struct {
	Config conf.Config
	Logger logg.Logger
}

func NewServer(config conf.Config, logger logg.Logger) *Serv {
	return &Serv{
		Config: config,
		Logger: logger,
	}
}

func (s *Serv) Run(router Router, mv middleware.Mdlwarer) {
	s.Logger.RaiseFatal(
		"server is not running",
		http.ListenAndServe(s.Config.GetRunAddress(), router.Run(mv)))
}
