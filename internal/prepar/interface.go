package prepar

import "net/http"

type ResPreper interface {
	OkWithJsonAndCookie(w http.ResponseWriter, data []byte, cookie *http.Cookie)
	BadOrConflWithErrJson(w http.ResponseWriter, status int, err error)
	UnauthorizedWithJson(w http.ResponseWriter, mess []byte)
}
