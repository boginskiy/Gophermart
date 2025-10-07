package middleware

import "net/http"

type Mdlwarer interface {
	WithAuth(next http.Handler) http.Handler
}
