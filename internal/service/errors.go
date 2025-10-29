package service

import "fmt"

type ErrHTTP struct {
	s string
	c int
}

func NewErrHTTP(code int) *ErrHTTP {
	return &ErrHTTP{
		s: fmt.Sprintf(`{"mess": "HTTP request failed with status %d"}`, code),
		c: code,
	}
}

func (e *ErrHTTP) Error() string {
	return e.s
}

func (e *ErrHTTP) TakeCode() int {
	return e.c
}
