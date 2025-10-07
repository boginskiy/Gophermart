package server

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/middleware"
)

type Router interface {
	Run(middleware.Mdlwarer) http.Handler
}

type Server interface {
	Run() error
}
