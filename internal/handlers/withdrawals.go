package handlers

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/service"
	"github.com/go-chi/chi"
)

type WithdrawalsHandlers struct {
	WithdrawalServ service.WithdrawalSrvcer
	ResPrep        prepar.ResPreper
}

func NewWithdrawalsHandlers(withdrawalSrvcer service.WithdrawalSrvcer, resPreper prepar.ResPreper) *WithdrawalsHandlers {
	return &WithdrawalsHandlers{
		WithdrawalServ: withdrawalSrvcer, ResPrep: resPreper,
	}
}

func (wh *WithdrawalsHandlers) RegisterRoutes(r chi.Router) {
	r.Get("/withdrawals", wh.GetWithdrawalHistory)
}

func (wh *WithdrawalsHandlers) GetWithdrawalHistory(w http.ResponseWriter, r *http.Request) {
	dataByte, err := wh.WithdrawalServ.GetWithdrawalHistory(r)

	// Ошибка сервера
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Исторические данные отсутствуют
	if 0 == len(dataByte) {
		wh.ResPrep.ResWithJson(w, service.MessNoOrders, http.StatusNoContent)
		return
	}

	wh.ResPrep.ResWithJson(w, dataByte, http.StatusOK)
}
