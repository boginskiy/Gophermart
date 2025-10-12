package prepar

import (
	"net/http"
)

type ResPrep struct {
}

func NewResPrep() *ResPrep {
	return &ResPrep{}
}

func (r *ResPrep) ResWithJson(
	w http.ResponseWriter,
	data []byte,
	status int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (r *ResPrep) ResWithJsonAndCookie(
	w http.ResponseWriter,
	data []byte,
	cookie *http.Cookie,
	status int) {

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
