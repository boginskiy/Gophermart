package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

type OrdersHandlers struct {
}

func NewOrdersHandlers() *OrdersHandlers {
	return &OrdersHandlers{}
}

func (oh *OrdersHandlers) RegisterRoutes(r chi.Router) {
	r.Post("/orders", oh.UploadOrderNumber)
	r.Get("/orders", oh.GetUploadedOrders)
}

func (oh *OrdersHandlers) UploadOrderNumber(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UploadOrderNumber"))
}

func (oh *OrdersHandlers) GetUploadedOrders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetUploadedOrders"))
}
