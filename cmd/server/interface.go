package server

import "net/http"

type Router interface {
	Run() http.Handler
}

type Server interface {
	Run() error
}
