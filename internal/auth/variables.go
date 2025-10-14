package auth

import "errors"

// Errors
var (
	ErrLoginPasswordIsBad  = errors.New(`{"error": "incorrect data for registration"}`)
	ErrLoginPasswordIsBad2 = errors.New(`{"error": "incorrect data for authorization"}`)
	ErrLogindNotFound      = errors.New(`{"error": "login not found"}`)
	ErrPasswordNotValid    = errors.New(`{"error": "password not valid"}`)
	ErrLoginNotUnic        = errors.New(`{"error": "login not unic, try again"}`)
	ErrCreateUser          = errors.New(`{"error": "user has not been created"}`)

	ErrTokenNotValid = errors.New(`{"error": "token not valid"}`)
	ErrTokenIsBad    = errors.New(`{"error": "token is bad"}`)
)

// Message
var (
	MessWelcome       = []byte(`{"mess": "welcome to Gophermart!"}`)
	MessNeedRegOrAuth = []byte(`{"mess": "go registration or authorization"}`)
)

// Ctx
type CtxKey struct{}

var CtxUserLogin = CtxKey{}
var CtxUserRole = CtxKey{}
var CtxUserID = CtxKey{}
