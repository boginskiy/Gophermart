package auth

import "net/http"

type Auther interface {
	Registration(req *http.Request) ([]byte, *http.Cookie, error)
	Authentication(req *http.Request) ([]byte, error)
}

type JWTokener interface {
	CheckOfValidToken(fullToken string) (login, role string, err error)
	CreateToken(login, role string) (token string, err error)
}
