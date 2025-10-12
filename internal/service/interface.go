package service

import "net/http"

type Srvcer interface {
}

type OrderSrvcer interface {
	UploadOrder(r *http.Request) ([]byte, error)
}

type OrderChecker interface {
	CheckDigits(digits string) bool
	GenDigits(long int) string
}
