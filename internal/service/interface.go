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
	GetOrders(r *http.Request) ([]byte, error)
}

type OrderChecker interface {
	CheckDigits(digits string) bool
	GenDigits(long int) string
}

type GatewaySrvcer interface {
	ConsumerAccruals(ctx context.Context)
	ConsumerOrders(ctx context.Context)
}

type BalanceSrvcer interface {
	GetBalance(req *http.Request) ([]byte, error)
	GetWithdrawal(req *http.Request) ([]byte, error)
}

type WithdrawalSrvcer interface {
	GetWithdrawalHistory(req *http.Request) ([]byte, error)
}
