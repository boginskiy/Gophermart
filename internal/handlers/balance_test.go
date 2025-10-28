package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boginskiy/Gophermart/cmd/client"
	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/cmd/server"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/boginskiy/Gophermart/internal/tests"
	"github.com/boginskiy/Gophermart/internal/tests/repotest"
	"github.com/boginskiy/Gophermart/pkg"
	"github.com/stretchr/testify/assert"
)

func TestBalanceHandlers(t *testing.T) {
	// Start server
	handler, _ := InitializationBalance()
	serv := httptest.NewServer(handler)
	defer serv.Close()

	// Client
	TUser := client.NewTUser(1, false)

	// Check Authentication
	testUserAuthentication(t, serv, TUser)
	// Check CurrentBalance
	testGetCurrentBalance(t, serv, TUser)
	// Check RequestWithdrawal
	testRequestWithdrawal(t, serv, TUser)
}

func InitializationBalance() (handler http.Handler, config *conf.ArgsENV) {
	// Инициализация
	businessLog := tests.NewTestLogg()
	infraLog := tests.NewTestLogg()
	appLog := tests.NewTestLogg()
	storeDB := tests.NewTestDB()

	config = conf.NewArgsENV(appLog)

	// Defer
	defer businessLog.Close()
	defer infraLog.Close()
	defer storeDB.Close()
	defer appLog.Close()

	// Repo
	repoLoyaltyOrders := repotest.NewTestRepoLoyaltyOrders(storeDB)
	repoOrders := repotest.NewTestRepoOrders(storeDB)
	repoUsers := repotest.NewTestRepoUsers(storeDB)

	// Services
	JWTServ := auth.NewJWTServ(config, appLog)
	auth := auth.NewAuth(config, appLog, repoUsers, JWTServ)
	orderChecker := pkg.NewLuna()
	balanceSrv := service.NewBalanceServ(
		config, businessLog, repoOrders, repoLoyaltyOrders, orderChecker)

	// Rsponse
	resPrep := prepar.NewResPrep()
	// Caps
	CAP := tests.NewTestHandlers()

	// Middleware
	mdlWare := middleware.NewMiddleware(config, appLog, auth, resPrep)
	// Handler
	balanceHdlrs := handlers.NewBalanceHandlers(balanceSrv, resPrep)
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)
	// Router
	router := server.NewRoute(authHdlrs, CAP, balanceHdlrs, CAP)
	return router.Run(mdlWare), config
}

func testUserAuthentication(t *testing.T, serv *httptest.Server, client *client.TUser) {
	tests := []struct {
		name          string
		bodyReq       string
		contentRes    string
		statusCodeRes int
	}{
		{name: "send bad login user's data", bodyReq: `{"LOGIN":"TestMan", "password":"123456"}`, contentRes: "application/json", statusCodeRes: 400},
		{name: "reqistration data not valid", bodyReq: `{"login":"TestMan", "password":"$$$$$$"}`, contentRes: "application/json", statusCodeRes: 401},
		{name: "send normal login user's data", bodyReq: `{"login":"TestMan", "password":"123456"}`, contentRes: "application/json", statusCodeRes: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			url := serv.URL + "/api/user/login"
			method := "POST"
			req, err := http.NewRequest(method, url, strings.NewReader(tt.bodyReq))
			assert.NoError(t, err)

			// Client
			res, err := client.C.Do(req)
			assert.NoError(t, err)

			// Check
			assert.Equal(t, tt.statusCodeRes, res.StatusCode)
			assert.Equal(t, tt.contentRes, res.Header.Get("Content-Type"))
		})
	}
}

func testGetCurrentBalance(t *testing.T, serv *httptest.Server, client *client.TUser) {
	url := serv.URL + "/api/user/balance"
	method := "GET"
	statusCodeRes := 200

	req, err := http.NewRequest(method, url, http.NoBody)
	assert.NoError(t, err)
	res, err := client.C.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, statusCodeRes, res.StatusCode)
}

func testRequestWithdrawal(t *testing.T, serv *httptest.Server, client *client.TUser) {
	tests := []struct {
		name          string
		bodyReq       string
		statusCodeRes int
		contentRes    string
	}{
		{name: "appl form without data", bodyReq: `{}`, statusCodeRes: 422, contentRes: "application/json"},
		{name: "appl form with bad order", bodyReq: `{"order": "12345678","sum": 0}`, statusCodeRes: 422, contentRes: "application/json"},
		{name: "appl form is ok", bodyReq: `{"order": "68379148","sum": 100}`, statusCodeRes: 200, contentRes: "application/json"},
		{name: "appl form is ok", bodyReq: `{"order": "68379148","sum": 1}`, statusCodeRes: 402, contentRes: "application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			url := serv.URL + "/api/user/balance/withdraw"
			method := "POST"
			req, err := http.NewRequest(method, url, strings.NewReader(tt.bodyReq))
			assert.NoError(t, err)

			// Response
			res, err := client.C.Do(req)
			assert.NoError(t, err)

			// Check
			assert.Equal(t, tt.statusCodeRes, res.StatusCode)
			assert.Equal(t, tt.contentRes, res.Header.Get("Content-Type"))
		})
	}
}
