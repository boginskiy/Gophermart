package service

import "errors"

var (
	ErrLoginPasswordIsBad = errors.New(`{"mess": "incorrect registration data"}`)
	ErrLoginNotUnic       = errors.New(`{"mess": "login is not unic"}`)
	ErrCreateUser         = errors.New(`{"mess": "user has not been created"}`)
)
