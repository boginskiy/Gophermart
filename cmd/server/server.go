package server

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/logg"
)

type Serv struct {
	args   string
	Logger logg.Logger
}

func NewServer(args string, logger logg.Logger) *Serv {
	return &Serv{
		args:   args,
		Logger: logger,
	}
}

func (s *Serv) Run(router Router) {
	s.Logger.RaiseFatal(
		"server is not running",
		http.ListenAndServe(`:8080`, router.Run()))
}
