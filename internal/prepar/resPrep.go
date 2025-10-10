package prepar

import (
	"net/http"
)

type ResPrep struct {
}

func NewResPrep() *ResPrep {
	return &ResPrep{}
}

func (r *ResPrep) UnauthorizedWithJson(w http.ResponseWriter, mess []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(mess)
}

func (r *ResPrep) BadOrConflWithErrJson(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func (r *ResPrep) OkWithJsonAndCookie(w http.ResponseWriter, data []byte, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
