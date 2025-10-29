package handlers

import (
	"net/http"

	"github.com/boginskiy/Gophermart/internal/prepar"
	"github.com/boginskiy/Gophermart/internal/service"

	"github.com/go-chi/chi"
)

type OrdersHandlers struct {
	OrderServ service.OrderSrvcer
	ResPrep   prepar.ResPreper
}

func NewOrdersHandlers(orderSrvcer service.OrderSrvcer, resPreper prepar.ResPreper) *OrdersHandlers {
	return &OrdersHandlers{OrderServ: orderSrvcer, ResPrep: resPreper}
}

func (oh *OrdersHandlers) RegisterRoutes(r chi.Router) {
	r.Post("/orders", oh.UploadOrderNumber)
	r.Get("/orders", oh.GetUploadedOrders)
}

func (oh *OrdersHandlers) UploadOrderNumber(w http.ResponseWriter, r *http.Request) {
	dataByte, err := oh.OrderServ.UploadOrder(r)

	// Невалидный номер заказа
	if err == service.ErrOrderNumber {
		oh.ResPrep.ResWithJSON(w, []byte(err.Error()), http.StatusUnprocessableEntity)
		return
	}

	// Дублирование заявки
	if err == service.ErrRepeatOrder {
		oh.ResPrep.ResWithJSON(w, []byte(err.Error()), http.StatusOK)
		return
	}

	// Заявка принадлежит другому пользователю
	if err == service.ErrAlienOrder {
		oh.ResPrep.ResWithJSON(w, []byte(err.Error()), http.StatusConflict)
		return
	}

	// Другие ошибки
	if err != nil {
		oh.ResPrep.ResWithJSON(w, []byte(err.Error()), http.StatusBadRequest)
		return
	}

	oh.ResPrep.ResWithJSON(w, dataByte, http.StatusAccepted)
}

func (oh *OrdersHandlers) GetUploadedOrders(w http.ResponseWriter, r *http.Request) {
	dataByte, err := oh.OrderServ.GetOrders(r)

	// Ошибки сервера
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Заказы отсутствуют
	if len(dataByte) == 0 {
		oh.ResPrep.ResWithJSON(w, service.MessNoOrders, http.StatusNoContent)
		return
	}

	oh.ResPrep.ResWithJSON(w, dataByte, http.StatusOK)
}
