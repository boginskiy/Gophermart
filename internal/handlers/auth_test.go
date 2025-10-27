package handlers_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/cmd/server"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/handlers"
	"github.com/boginskiy/Gophermart/internal/middleware"
	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/internal/tests"
	"github.com/boginskiy/Gophermart/internal/tests/repotest"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers2(t *testing.T) {
	// Инициализация
	appLog := tests.NewTestLogg()
	config := conf.NewArgsENV(appLog)
	infraLog := tests.NewTestLogg()
	businessLog := tests.NewTestLogg()

	storeDB := tests.NewTestDB()

	// Defer
	defer businessLog.Close()
	defer infraLog.Close()
	defer storeDB.Close()
	defer appLog.Close()

	repoUsers := repository.NewRepoUsers(config, infraLog, storeDB)
	JWTServ := auth.NewJWTServ(config, appLog)

	auth := auth.NewAuth(config, appLog, repoUsers, JWTServ)
	resPrep := prepar.NewResPrep()

	// Middleware
	mdlWare := middleware.NewMiddleware(config, appLog, auth, resPrep)
	// Handler
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)
	// Router
	router := server.NewRoute(authHdlrs, tests.NewTestHandlers(), tests.NewTestHandlers(), tests.NewTestHandlers())

	router.Run(mdlWare)

	// http.ListenAndServe(s.config.GetRunAddress(), router.Run(mv))

}

func TestAythHandler(t *testing.T) {
	// Handler
	auth, resPrep := InitializationAuth()
	authHdlrs := handlers.NewAuthHandlers(auth, resPrep)

	// Testing
	testRegisterUser(t, authHdlrs)
	testLoginUser(t, authHdlrs)
}

func InitializationAuth() (*auth.Auth, *prepar.ResPrep) {
	appLog := tests.NewTestLogg()
	storeDB := tests.NewTestDB()
	config := conf.NewArgsENV(appLog)

	repoUsers := repotest.NewTestRepoUsers(storeDB)
	JWTServ := auth.NewJWTServ(config, appLog)
	auth := auth.NewAuth(config, appLog, repoUsers, JWTServ)

	resPrep := prepar.NewResPrep()
	return auth, resPrep
}

func testRegisterUser(t *testing.T, handler *handlers.AuthHandlers) {

	type tRes struct {
		code        int
		body        string
		contentType string
	}

	type tReq struct {
		url    string
		method string
		body   string
	}

	tests := []struct {
		name string
		tRes tRes
		tReq tReq
	}{
		{
			name: "send bad register user's data",
			tRes: tRes{
				code:        400,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/register",
				method: "POST",
				body:   `{"LOGIN":"Vasia", "PASSWORD":"1234"}`},
		},
		{
			name: "send normal register user's data",
			tRes: tRes{
				code:        200,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/register",
				method: "POST",
				body:   `{"login":"Vasia", "password":"1234"}`},
		},
		{
			name: "repeate send normal register user's data",
			tRes: tRes{
				code:        409,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/register",
				method: "POST",
				body:   `{"login":"Vasia", "password":"1234"}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			body := bytes.NewReader([]byte(tt.tReq.body))
			request := httptest.NewRequest(tt.tReq.method, tt.tReq.url, body)

			// Recorder ~ http.ResponseWriter
			w := httptest.NewRecorder()

			// Handler
			handler.RegisterUser(w, request)

			// Response
			res := w.Result()

			// Check
			assert.Equal(t, res.Header.Get("Content-Type"), tt.tRes.contentType)
			assert.Equal(t, res.StatusCode, tt.tRes.code)

		})
	}
}

func testLoginUser(t *testing.T, handler *handlers.AuthHandlers) {
	type tRes struct {
		code        int
		body        string
		contentType string
	}

	type tReq struct {
		url    string
		method string
		body   string
	}

	tests := []struct {
		name string
		tRes tRes
		tReq tReq
	}{
		{
			name: "send bad login user's data",
			tRes: tRes{
				code:        400,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/login",
				method: "POST",
				body:   `{"LOGIN":"Vasia", "PASSWORD":"1234"}`},
		},
		{
			name: "send normal login user's data",
			tRes: tRes{
				code:        200,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/login",
				method: "POST",
				body:   `{"login":"TestMan", "password":"123456"}`},
		},
		{
			name: "user doesn't have registration",
			tRes: tRes{
				code:        401,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/login",
				method: "POST",
				body:   `{"login":"Nemo", "password":"&&&&"}`},
		},

		{
			name: "reqistration data not valid ",
			tRes: tRes{
				code:        401,
				contentType: "application/json"},

			tReq: tReq{
				url:    "/login",
				method: "POST",
				body:   `{"login":"TestMan", "password":"$$$$$$"}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			body := bytes.NewReader([]byte(tt.tReq.body))
			req := httptest.NewRequest(tt.tReq.method, tt.tReq.url, body)

			w := httptest.NewRecorder()
			handler.LoginUser(w, req)

			// Response
			res := w.Result()

			// Check
			assert.Equal(t, res.Header.Get("Content-Type"), tt.tRes.contentType)
			assert.Equal(t, res.StatusCode, tt.tRes.code)
		})
	}
}
