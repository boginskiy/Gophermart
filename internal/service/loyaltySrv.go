package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	repo "github.com/boginskiy/Gophermart/internal/repository"
	repos "github.com/boginskiy/Gophermart/internal/repository"
	mod "github.com/boginskiy/Gophermart/models"
)

type Loyalty struct {
	Repo       repo.RepoOrdersTber
	ChAccruals chan *mod.Accrual
	ChOrders   chan *mod.Order
	Core       *CoreSrv

	ctx       context.Context
	semaphore chan struct{}
	host      string
	path      string
}

func NewLoyalty(ctx context.Context, chOrders chan *mod.Order, core *CoreSrv, repo repos.RepoOrdersTber) *Loyalty {
	item := &Loyalty{
		ChAccruals: make(chan *mod.Accrual, 10),
		ChOrders:   chOrders,
		Core:       core,
		Repo:       repo,

		ctx:       ctx,
		semaphore: make(chan struct{}, core.Args.GetMaxConcurrentReq()),
		host:      core.Args.GetACCRUAL_SYSTEM_ADDRESS(),
		path:      "/api/orders/", // TODO! Вывести на default args
	}

	// Запуск фоновой работы прокси модуля 'Loyalty'
	go item.ConsumerOrders(ctx)

	return item
}

// Сonsumer
func (l *Loyalty) ConsumerOrders(ctx context.Context) {
	defer close(l.ChAccruals)

	// Может тут сделать локальные контексты ?
	// В общем с этим надо разбираться!

	for {
		select {

		// Обработка заявок
		case order := <-l.ChOrders:
			l.SendOrder(order)

		// Завершение обработки
		case <-ctx.Done():
			return
		}
	}
}

func (l *Loyalty) fetchData(url string) (accrual *mod.Accrual) {
	// Сетевой запрос в 'bonusService'
	res, err := http.Get(url)
	if err != nil {
		l.Core.Logg.RaiseError("Loyalty>fetchData>http.Get", err)
		return accrual
	}

	// Обработка кодов ответа: 204, 429, 500 ...
	if res.StatusCode != http.StatusOK {
		l.Core.Logg.RaiseInfo(fmt.Sprintf("bonusService answered %d status", res.StatusCode))
		return accrual
	}

	// Читаем тело ответа от 'bonusService'
	dataByte, err := io.ReadAll(res.Body)
	if err != nil {
		l.Core.Logg.RaiseError("Loyalty>fetchData>ReadAll", err)
		return accrual
	}
	defer res.Body.Close()

	// Сериализуем
	err = json.Unmarshal(dataByte, accrual)
	if err != nil {
		l.Core.Logg.RaiseError("Loyalty>fetchData>Unmarshal", err)
	}
	return accrual
}

func (l *Loyalty) SendOrder(order *mod.Order) {
	url := fmt.Sprintf("%s%s%s", l.host, l.path, order.Code)

	go func(ctx context.Context, url string) {

		l.semaphore <- struct{}{}

		// Получаем данные
		if v := l.fetchData(url); v != nil {
			l.ChAccruals <- v
		}

		<-l.semaphore

	}(l.ctx, url)

}

// Хер знает как сделать
