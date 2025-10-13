package service

import (
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/repository"
)

type OrderSrv struct {
	Args       config.Argser
	Logg       logg.Logger
	Repo       repository.RepoOrdersTber
	OrderCheck OrderChecker
}

func NewOrderSrv(
	argser config.Argser,
	logger logg.Logger,
	repo repository.RepoOrdersTber,
	orderChecker OrderChecker) *OrderSrv {

	return &OrderSrv{
		Args:       argser,
		Logg:       logger,
		Repo:       repo,
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
