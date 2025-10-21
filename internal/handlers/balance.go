package handlers

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/go-chi/chi"
)

type BalanceHandlers struct {
	BalanceServ service.BalanceSrvcer
	ResPrep     prepar.ResPreper
}

func NewBalanceHandlers(balanceSrvcer service.BalanceSrvcer, resPreper prepar.ResPreper) *BalanceHandlers {
	return &BalanceHandlers{
		BalanceServ: balanceSrvcer, ResPrep: resPreper}
}

func (bh *BalanceHandlers) RegisterRoutes(r chi.Router) {
	r.Get("/balance", bh.GetCurrentBalance)
	r.Post("/balance/withdraw", bh.RequestWithdrawal)
}

func (bh *BalanceHandlers) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	dataByte, err := bh.BalanceServ.GetBalance(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bh.ResPrep.ResWithJSON(w, dataByte, http.StatusOK)
}

func (bh *BalanceHandlers) RequestWithdrawal(w http.ResponseWriter, r *http.Request) {
	dataByte, err := bh.BalanceServ.GetWithdrawal(r)

	// Неверный номер заказа
	if err == service.ErrOrderNumber {
		bh.ResPrep.ResWithJSON(w, []byte(err.Error()), http.StatusUnprocessableEntity)
		return
	}

	// Ошибка превышения лимита бонуса
	if err == service.ErrBonuseLimit {
		bh.ResPrep.ResWithJSON(w, []byte(err.Error()), http.StatusPaymentRequired)
		return
	}

	// Ошибки сервера
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bh.ResPrep.ResWithJSON(w, dataByte, http.StatusOK)
}
