package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type BalanceHandlers struct{}

func NewBalanceHandlers() *BalanceHandlers {
	return &BalanceHandlers{}
}

func (bh *BalanceHandlers) RegisterRoutes(r chi.Router) {
	r.Get("/balance", bh.GetCurrentBalance)
	r.Post("/balance/withdraw", bh.RequestWithdrawal)
}

func (bh *BalanceHandlers) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetCurrentBalance"))
}

func (bh *BalanceHandlers) RequestWithdrawal(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("RequestWithdrawal"))
}
