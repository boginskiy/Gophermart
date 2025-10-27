package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/logg"
	repo "github.com/boginskiy/Gophermart/internal/repository"
	mod "github.com/boginskiy/Gophermart/models"
)

type OrderSrv struct {
	Chan       chan *mod.Order
	Config     conf.Config
	Logger     logg.Logger
	Repo       repo.RepoOrdersTber
	OrderCheck OrderChecker
}

func NewOrderSrv(
	ch chan *mod.Order,
	config conf.Config,
	logger logg.Logger,
	repoOrders repo.RepoOrdersTber,
	orderCheck OrderChecker) *OrderSrv {

	return &OrderSrv{
		Chan:       ch,
		Config:     config,
		Logger:     logger,
		Repo:       repoOrders,
		OrderCheck: orderCheck}
}

// sendOrdersToGateWay - метод для отправки заказов в сервис 'GateWay'
func (o *OrderSrv) sendOrdersToGateWay(order *mod.Order) []byte {
	select {
	case o.Chan <- order:
		return MessNewOrder
	default:
		// Если переполнение очереди заказов, сообщаем, что сервис перегружен
		return MessOverLoadOrder
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
	if !o.OrderCheck.CheckDigits(orderCode) {
		return nil, ErrOrderNumber
	}

	// Пробуем достать запись с orderCode
	userOrder, err := o.Repo.ReadWithUser(context.TODO(), orderCode)

	// Запись с таким orderCode уже есть в БД
	if err == nil {
		loginReq := req.Context().Value(auth.CtxUserLogin)
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
	userID := req.Context().Value(auth.CtxUserID)
	newOrder := mod.NewOrder(orderCode, userID.(int64))

	_, err = o.Repo.Create(context.TODO(), newOrder)
	if err != nil {
		return nil, err
	}

	// Передаем через общий канал 'ChOrders' заказ в gateway-модуль
	mess := o.sendOrdersToGateWay(newOrder)

	return mess, nil
}

func (o *OrderSrv) GetOrders(req *http.Request) ([]byte, error) {
	userID := req.Context().Value(auth.CtxUserID)
	orders, err := o.Repo.ReadOrdersWithSort(context.TODO(), userID.(int64))

	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return EmptySliceOfBytes, nil
	}

	return json.Marshal(orders)
}
