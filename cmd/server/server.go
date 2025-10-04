package server

import (
	"net/http"
)

type Serv struct {
	args string
	logg string
}

func NewServer(args string, logg string) *Serv {
	return &Serv{
		args: args,
		logg: logg,
	}
}

func (s *Serv) Run(router Router) error {
	return http.ListenAndServe(`:8080`, router.Run())
}
