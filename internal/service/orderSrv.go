package service

import (
	"context"
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	repo "github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/models"
	mod "github.com/boginskiy/Gophermart/models"
)

type OrderSrv struct {
	ChOrders chan *models.Order
	Core     *CoreSrv
	Repo     repo.RepoOrdersTber
}

func NewOrderSrv(chOrders chan *models.Order, c *CoreSrv, r repo.RepoOrdersTber) *OrderSrv {
	return &OrderSrv{
		ChOrders: chOrders,
		Core:     c,
		Repo:     r,
	}
}

func (o *OrderSrv) UploadOrder(req *http.Request) ([]byte, error) {
	// Вынимаем номер заказа
	dataByte, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}

	orderCode := string(dataByte)

	// Проверка номера заказа алгоритмом 'Луна'
	if !o.Core.OrderCheck.CheckDigits(orderCode) {
		return nil, ErrOrderNumber
	}

	// Пробуем достать запись с orderCode
	userOrder, err := o.Repo.ReadWithUser(context.TODO(), orderCode)

	// Запись с таким orderCode уже есть в БД
	if err == nil {
		loginReq := o.Core.takeParamFromAuth(req, auth.CtxUserLogin)
		loginDB := userOrder.Login

		if loginReq.(string) != loginDB {
			// Запись уже была загружена другим пользователем
			return nil, ErrAlienOrder
		} else {
			// Запись уже была загружена текущим пользователем
			return nil, ErrRepeatOrder
		}
	}

	// Создаем новую запись с заказом. Сохраняем в БД
	userID := o.Core.takeParamFromAuth(req, auth.CtxUserID)
	newOrder := mod.NewOrder(orderCode, userID.(int64))

	_, err = o.Repo.Create(context.TODO(), newOrder)
	if err != nil {
		return nil, err
	}

	// Передаем через общий канал 'newOrder' в модуль, который
	// взаимодействует с сервисом расчета баллов лояльности
	o.ChOrders <- newOrder

	return MessNewOrder, nil
}
