package service

import (
	"context"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
)

type CoreSrvcer interface {
	takeParamFromAuth(req *http.Request, p auth.CtxKey) any
}

type OrderSrvcer interface {
	UploadOrder(r *http.Request) ([]byte, error)
}

type OrderChecker interface {
	CheckDigits(digits string) bool
	GenDigits(long int) string
}

type GatewaySrvcer interface {
	ConsumerAccruals(ctx context.Context)
	ConsumerOrders(ctx context.Context)
}
