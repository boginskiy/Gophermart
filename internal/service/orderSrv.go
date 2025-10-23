package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	repo "github.com/boginskiy/Gophermart/internal/repository"
	mod "github.com/boginskiy/Gophermart/models"
)

type OrderSrv struct {
	ChOrders chan *mod.Order
	Core     *CoreSrv
	Repo     repo.RepoOrdersTber
}

func NewOrderSrv(chOrders chan *mod.Order, c *CoreSrv, r repo.RepoOrdersTber) *OrderSrv {
	return &OrderSrv{
		ChOrders: chOrders,
		Core:     c,
		Repo:     r,
	}
}

// sendOrdersToGateWay - метод для отправки заказов в сервис 'GateWay'
func (o *OrderSrv) sendOrdersToGateWay(order *mod.Order) []byte {
	log.Println("8>>", order)
	select {
	case o.ChOrders <- order:
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
	if !o.Core.OrderCheck.CheckDigits(orderCode) {
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
		log.Println(err)
		return nil, err
	}

	if len(orders) == 0 {
		return EmptySliceOfBytes, nil
	}

	return json.Marshal(orders)
}
