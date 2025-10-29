package prepar

import "net/http"

type ResPreper interface {
	ResWithJSONAndCookie(w http.ResponseWriter, data []byte, cookie *http.Cookie, status int)
	ResWithJSON(w http.ResponseWriter, data []byte, status int)
}
