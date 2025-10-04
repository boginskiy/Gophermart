package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type WithdrawalsHandlers struct{}

func NewWithdrawalsHandlers() *WithdrawalsHandlers {
	return &WithdrawalsHandlers{}
}

func (wh *WithdrawalsHandlers) RegisterRoutes(r chi.Router) {
	r.Get("/withdrawals", wh.GetWithdrawalHistory)
}

func (wh *WithdrawalsHandlers) GetWithdrawalHistory(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetWithdrawalHistory"))
}
