package service

import "net/http"

type Srvcer interface {
}

type UserSrvcer interface {
	Authentication(req *http.Request) ([]byte, error)
	Registration(req *http.Request) ([]byte, error)
}
