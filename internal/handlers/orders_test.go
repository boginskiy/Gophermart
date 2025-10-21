package handlers_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/boginskiy/Gophermart/internal/tests"
	"github.com/boginskiy/Gophermart/internal/tests/repotest"
	"github.com/boginskiy/Gophermart/models"
	"github.com/boginskiy/Gophermart/pkg"
	"github.com/stretchr/testify/assert"
)

func TestOrdersHandlers(t *testing.T) {
	// Chan
	chOrders := make(chan *models.Order, 10)
	defer close(chOrders)

	// Initialization
	orderSrv, resPrep := InitializationOrder(chOrders)
	orderHdlrs := handlers.NewOrdersHandlers(orderSrv, resPrep)

	// Testing
	testUploadOrderNumber(t, orderHdlrs)
	testGetUploadedOrders(t, orderHdlrs)
}

func InitializationOrder(chOrders chan *models.Order) (*service.OrderSrv, *prepar.ResPrep) {
	appLog := tests.NewTestLogg()
	storeDB := tests.NewTestDB()
	args := config.NewArgsENV(appLog)

	repoOrders := repotest.NewTestRepoOrders(storeDB)

	orderChecker := pkg.NewLuna()
	coreSrv := service.NewCoreSrv(args, appLog, orderChecker)
	orderSrv := service.NewOrderSrv(chOrders, coreSrv, repoOrders)

	resPrep := prepar.NewResPrep()

	return orderSrv, resPrep
}

func testUploadOrderNumber(t *testing.T, handler *handlers.OrdersHandlers) {
	type tRes struct {
		code        int
		body        string
		contentType string
	}

	type tReq struct {
		url      string
		method   string
		body     string
		ctxKey   auth.CtxKey
		ctxValue any
	}

	tests := []struct {
		name string
		tRes tRes
		tReq tReq
	}{
		{
			name: "number of order is not valid",
			tRes: tRes{
				code:        422,
				contentType: "application/json"},

			tReq: tReq{
				url:      "/orders",
				method:   "POST",
				body:     "12345678",
				ctxKey:   auth.CtxUserLogin,
				ctxValue: "TestMan"},
		},
		{
			name: "number of order is duplication",
			tRes: tRes{
				code:        200,
				contentType: "application/json"},

			tReq: tReq{
				url:      "/orders",
				method:   "POST",
				body:     "68379148",
				ctxKey:   auth.CtxUserLogin,
				ctxValue: "TestMan"},
		},
		{
			name: "application belongs to another user",
			tRes: tRes{
				code:        409,
				contentType: "application/json"},

			tReq: tReq{
				url:      "/orders",
				method:   "POST",
				body:     "68379148",
				ctxKey:   auth.CtxUserLogin,
				ctxValue: "TestMan2"},
		},
		{
			name: "new load of application",
			tRes: tRes{
				code:        202,
				contentType: "application/json"},

			tReq: tReq{
				url:      "/orders",
				method:   "POST",
				body:     "96357934",
				ctxKey:   auth.CtxUserID,
				ctxValue: int64(1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			body := bytes.NewReader([]byte(tt.tReq.body))
			req := httptest.NewRequest(tt.tReq.method, tt.tReq.url, body)

			// Context with 'TestMan'
			ctx := context.WithValue(req.Context(), tt.tReq.ctxKey, tt.tReq.ctxValue)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.UploadOrderNumber(w, req)

			// Response
			res := w.Result()

			// Check
			assert.Equal(t, res.Header.Get("Content-Type"), tt.tRes.contentType)
			assert.Equal(t, res.StatusCode, tt.tRes.code)
		})
	}
}

func testGetUploadedOrders(t *testing.T, handler *handlers.OrdersHandlers) {
	type tRes struct {
		code        int
		body        string
		contentType string
	}

	type tReq struct {
		url      string
		method   string
		body     string
		ctxKey   auth.CtxKey
		ctxValue any
	}

	tests := []struct {
		name string
		tRes tRes
		tReq tReq
	}{
		{
			name: "there are no orders",
			tRes: tRes{
				code:        204,
				contentType: "application/json"},

			tReq: tReq{
				url:      "/orders",
				method:   "GET",
				body:     "96357934",
				ctxKey:   auth.CtxUserID,
				ctxValue: int64(2)},
		},
		{
			name: "there are orders",
			tRes: tRes{
				code:        200,
				contentType: "application/json"},

			tReq: tReq{
				url:      "/orders",
				method:   "GET",
				body:     "96357934",
				ctxKey:   auth.CtxUserID,
				ctxValue: int64(1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			req := httptest.NewRequest(tt.tReq.method, tt.tReq.url, nil)

			// Context
			ctx := context.WithValue(req.Context(), tt.tReq.ctxKey, tt.tReq.ctxValue)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.GetUploadedOrders(w, req)

			// Response
			res := w.Result()

			// Check
			assert.Equal(t, res.Header.Get("Content-Type"), tt.tRes.contentType)
			assert.Equal(t, res.StatusCode, tt.tRes.code)
		})
	}
}
