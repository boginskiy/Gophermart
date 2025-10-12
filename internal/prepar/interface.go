package prepar

import "net/http"

type ResPreper interface {
	ResWithJsonAndCookie(w http.ResponseWriter, data []byte, cookie *http.Cookie, status int)
	ResWithJson(w http.ResponseWriter, data []byte, status int)
}
