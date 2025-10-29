package store

import "errors"

var (
	ErrOpeningDB  = errors.New(`{"mess": "database opening error"}`)
	ErrCreatingTb = errors.New(`{"mess": "tables creating error"}`)
)
