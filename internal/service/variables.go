package service

import "errors"

// Errors
var (
	ErrLoginPasswordIsBad = errors.New(`{"mess": "incorrect registration data"}`)
	ErrLoginNotUnic       = errors.New(`{"mess": "login is not unic"}`)
	ErrCreateUser         = errors.New(`{"mess": "user has not been created"}`)
	ErrOrderNumber        = errors.New(`{"mess": "number of order is not valid"}`)
	ErrAlienOrder         = errors.New(`{"mess": "application does not belong to you"}`)
	ErrRepeatOrder        = errors.New(`{"mess": "application is processing"}`)
)

// Mess
var MessNewOrder = []byte(`{"mess": "new order has been accepted for processing"}`)

// Empty
var EmptySliceOfBytes = []byte{}
