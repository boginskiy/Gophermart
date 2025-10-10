package auth

import "net/http"

type Auther interface {
	CheckAuthReq(req *http.Request) bool
	CheckToken(token string) (login, role string, err error)
	Registration(req *http.Request) ([]byte, *http.Cookie, error)
	Authentication(req *http.Request) ([]byte, *http.Cookie, error)
}

type JWTokener interface {
	CheckOfValidToken(fullToken string) (login, role string, err error)
	CreateToken(login, role string) (token string, err error)
}
