package repository

import "errors"

var (
	ErrType = errors.New(`{"mess": "incorrect type of item for db"}`)
)
