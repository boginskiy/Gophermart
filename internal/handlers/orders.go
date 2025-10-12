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
		oh.ResPrep.ResWithJson(w, []byte(err.Error()), http.StatusUnprocessableEntity)
		return
	}

	// Другие ошибки
	if err != nil {
		oh.ResPrep.ResWithJson(w, []byte(err.Error()), http.StatusBadRequest)
		return
	}

	oh.ResPrep.ResWithJson(w, dataByte, http.StatusOK)

}

func (oh *OrdersHandlers) GetUploadedOrders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetUploadedOrders"))
}
