package service

import (
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/repository"
)

type OrderSrv struct {
	Repo       repository.RepoDBer
	Logger     logg.Logger
	OrderCheck OrderChecker
}

func NewOrderSrv(repo repository.RepoDBer, logger logg.Logger, orderChecker OrderChecker) *OrderSrv {
	return &OrderSrv{
		Repo:       repo,
		Logger:     logger,
		OrderCheck: orderChecker,
	}
}

func (o *OrderSrv) UploadOrder(req *http.Request) ([]byte, error) {
	// Вынимаем номер заказа
	numOrder, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}

	// Проверка номера заказа алгоритмом 'Луна'
	if !o.OrderCheck.CheckDigits(string(numOrder)) {
		// Неверный формат номера заказа
		return nil, ErrOrderNumber
	}

	// Номер заказа
	// numOrder

	return nil, nil
}
