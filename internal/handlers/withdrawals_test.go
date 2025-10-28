package handlers_test

import (
	"encoding/json"
	"io"
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
	"github.com/boginskiy/Gophermart/models"
	"github.com/stretchr/testify/assert"
)

func InitializationWithdrawals() (handler http.Handler, config *conf.ArgsENV) {
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
	repoUsers := repotest.NewTestRepoUsers(storeDB)

	// Services
	JWTServ := auth.NewJWTServ(config, appLog)
	auth := auth.NewAuth(config, appLog, repoUsers, JWTServ)
	withdrawalsSrv := service.NewWithdrawalsSrv(config, businessLog, repoLoyaltyOrders)

	// Rsponse
	resPrep := prepar.NewResPrep()
	// Caps
	CAP := tests.NewTestHandlers()

	// Middleware
	mdlWare := middleware.NewMiddleware(config, appLog, auth, resPrep)
	// Handler
	withdrawalsHdlrs := handlers.NewWithdrawalsHandlers(withdrawalsSrv, resPrep)
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)
	// Router
	router := server.NewRoute(authHdlrs, CAP, CAP, withdrawalsHdlrs)
	return router.Run(mdlWare), config
}

func authentication(t *testing.T, serv *httptest.Server, client1, client2 *client.TUser) {
	tests := []struct {
		name          string
		client        *client.TUser
		bodyReq       string
		contentRes    string
		statusCodeRes int
	}{
		{name: "send normal login user's data",
			bodyReq:       `{"login":"TestMan", "password":"123456"}`,
			contentRes:    "application/json",
			statusCodeRes: 200,
			client:        client1},

		{name: "send normal login user2's data",
			bodyReq:       `{"login":"TestMan2", "password":"7890123"}`,
			contentRes:    "application/json",
			statusCodeRes: 200,
			client:        client2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			url := serv.URL + "/api/user/login"
			method := "POST"
			req, err := http.NewRequest(method, url, strings.NewReader(tt.bodyReq))
			assert.NoError(t, err)

			// Client
			res, err := tt.client.C.Do(req)
			assert.NoError(t, err)

			// Check
			assert.Equal(t, tt.statusCodeRes, res.StatusCode)
			assert.Equal(t, tt.contentRes, res.Header.Get("Content-Type"))
		})
	}
}

func testGetWithdrawalHistory(t *testing.T, serv *httptest.Server, client1, client2 *client.TUser) {
	tests := []struct {
		name            string
		client          *client.TUser
		statusCodeRes   int
		lengthOfBodyRes int
	}{
		{name: "client with some data", client: client1, statusCodeRes: 200, lengthOfBodyRes: 2},
		{name: "client without some data", client: client2, statusCodeRes: 204, lengthOfBodyRes: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			url := serv.URL + "/api/user/withdrawals"
			method := "GET"
			req, err := http.NewRequest(method, url, http.NoBody)
			assert.NoError(t, err)

			// Response
			res, err := tt.client.C.Do(req)
			assert.NoError(t, err)

			// Check
			assert.Equal(t, tt.statusCodeRes, res.StatusCode)

			// Check response body
			sl := []*models.LoyaltyOrder{}
			dataByte, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			json.Unmarshal(dataByte, &sl)
			assert.Equal(t, tt.lengthOfBodyRes, len(sl))
		})
	}
}

func TestWithdrawalsHandlers(t *testing.T) {
	// Start server
	handler, _ := InitializationWithdrawals()
	serv := httptest.NewServer(handler)
	defer serv.Close()

	// Clients
	TUser1 := client.NewTUser(1, false)
	TUser2 := client.NewTUser(2, false)

	// Authentication
	authentication(t, serv, TUser1, TUser2)
	// Check GetWithdrawalHistory
	testGetWithdrawalHistory(t, serv, TUser1, TUser2)
}
