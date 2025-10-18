package service

import (
	"errors"
)

const SIZE = 10

// Errors
var (
	ErrLoginPasswordIsBad = errors.New(`{"mess": "incorrect registration data"}`)
	ErrLoginNotUnic       = errors.New(`{"mess": "login is not unic"}`)
	ErrCreateUser         = errors.New(`{"mess": "user has not been created"}`)
	ErrOrderNumber        = errors.New(`{"mess": "number of order is not valid"}`)
	ErrAlienOrder         = errors.New(`{"mess": "application does not belong to you"}`)
	ErrRepeatOrder        = errors.New(`{"mess": "application is processing"}`)
	ErrBonuseLimit        = errors.New(`{"mess": "request bonuse is not valid"}`)
)

// Mess
var MessNewOrder = []byte(`{"mess": "new order has been accepted for processing"}`)
var MessOverLoadOrder = []byte(`{"mess": "service is overload, try request again"}`)
var MessNoOrders = []byte(`{"mess": "no orders"}`)

// Empty
var EmptySliceOfBytes = []byte{}
