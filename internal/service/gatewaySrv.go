package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	repo "github.com/boginskiy/Gophermart/internal/repository"
	repos "github.com/boginskiy/Gophermart/internal/repository"
	mod "github.com/boginskiy/Gophermart/models"
)

type GatewaySrv struct {
	Repo     repo.RepoOrdersTber
	ChOrders chan *mod.Order
	Core     *CoreSrv
	Ctx      context.Context

	chAccruals chan *mod.Accrual
	semaphore  chan struct{}
	host       string
	path       string
	timeWaite  time.Duration
	client     *http.Client
}

func NewGatewaySrv(ctx context.Context, chOrders chan *mod.Order, core *CoreSrv, repo repos.RepoOrdersTber) *GatewaySrv {
	item := &GatewaySrv{
		ChOrders: chOrders,
		Core:     core,
		Repo:     repo,
		Ctx:      ctx,

		semaphore:  make(chan struct{}, core.Args.GetMaxConcurrentReq()),
		timeWaite:  core.Args.TimeWaitingResponse(),
		host:       core.Args.GetSystemAddress(),
		chAccruals: make(chan *mod.Accrual, SIZE),
		path:       core.Args.GetSystemPath(),
		client:     &http.Client{},
	}

	// Запуск фоновой работы прокси модуля 'GatewaySrv'
	go item.ConsumerOrders(ctx)
	go item.ConsumerAccruals(ctx)

	return item
}

// ConsumerAccruals -
func (l *GatewaySrv) ConsumerAccruals(ctx context.Context) {
	// Каждые N-секунд осуществляем обновление статусов заказов в БД
	ticker := time.NewTicker(l.Core.Args.GetTimeTicker())
	// Читаем из канала. При закрытии канала данные должны упасть в БД
	slAccruals := make([]*mod.Accrual, 0, 10)

	select {

	// Данные поступили. Делаем запись в БД раз в N-сек
	case <-ticker.C:
		slAccruals = l.SendAccrualsToDB(slAccruals)

	// Собираем данные для последующей записи
	case accrual := <-l.chAccruals:
		slAccruals = append(slAccruals, accrual)

	case <-ctx.Done():
		// Дописываем остатки данных
		l.SendAccrualsToDB(slAccruals)
		return
	}
}

func (l *GatewaySrv) SendAccrualsToDB(accruals []*mod.Accrual) []*mod.Accrual {
	if 0 < len(accruals) {
		err := l.Repo.UpdateSetStatuses(context.TODO(), accruals)
		if err != nil {
			l.Core.Logg.RaiseInfo("GatewaySrv>SendAccrualsToDB: query is bad")
		} else {
			return accruals[:0]
		}
	}
	return accruals
}

// ConsumerOrders -
func (l *GatewaySrv) ConsumerOrders(ctx context.Context) {
	ticker := time.NewTicker(l.Core.Args.GetTimeTicker())
	slOrders := make([]*mod.Order, 0, 10)
	defer close(l.chAccruals)

	for {
		select {

		// Каждые N-секунд передаем новые заказы в обработку
		case <-ticker.C:
			slOrders = l.SendOrdersToDistantSrv(slOrders)

		// Обработка заявок
		case order := <-l.ChOrders:
			slOrders = append(slOrders, order)

		// Данные, которые останутся в канале или временном хранилище
		// Останутся в БД со статусом "NEW"
		case <-ctx.Done():
			return
		}
	}
}

func (l *GatewaySrv) fetchData(ctx context.Context, url string) (*mod.Accrual, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Обработка кодов ответа: 204, 429, 500
	if res.StatusCode != http.StatusOK {
		return nil, NewErrHTTP(res.StatusCode)
	}

	var accrual mod.Accrual
	err = json.NewDecoder(res.Body).Decode(&accrual)
	if err != nil {
		return nil, err
	}
	return &accrual, nil
}

func (l *GatewaySrv) procesData(order *mod.Order, accrual *mod.Accrual, err error) {
	// Логика обработки кодов StatusCode: 204, 429, 500
	if _, ok := err.(*ErrHTTP); !ok {

		// Если удаленный сервис по каким либо причинам, а именно:
		//    StatusCode 204 — заказ не зарегистрирован в системе расчёта
		//    StatusCode 429 — превышено количество запросов к сервису
		//    StatusCode 500 — внутренняя ошибка сервера
		// то передаем заявку в очередь на повторный запрос
		l.Core.Logg.RaiseError("GatewaySrv>sendOrderToDistantSrv>procesData", err)
		l.ChOrders <- order
		return
	}

	// Внутренние ошибки обработки
	if err != nil || accrual == nil {
		// Передаем заявку повторно в очередь обработки
		l.Core.Logg.RaiseError("GatewaySrv>sendOrderToDistantSrv>procesData", err)
		l.ChOrders <- order
		return
	}

	// Заявки со статусами 'REGISTERED', 'PROCESSING' не являются окончательными
	// их отправляем на повторную обработку
	if accrual.Status == "REGISTERED" || accrual.Status == "PROCESSING" {
		l.ChOrders <- order
		return
	}

	// Заявки со статусами 'INVALID', 'PROCESSED' являются окончательными
	// их отправляем далее для обновления статуса в БД
	if accrual.Status == "INVALID" || accrual.Status == "PROCESSED" {
		l.chAccruals <- accrual
		return
	}
}

func (l *GatewaySrv) sendOrderToDistantSrv(order *mod.Order, url string) {
	// Контекст с таймаутом + Semaphore
	ctx, cancel := context.WithTimeout(l.Ctx, l.timeWaite)
	defer func() { <-l.semaphore }()
	defer cancel()

	// Отправляем запрос на получение данных о бонусах
	accrual, err := l.fetchData(ctx, url)

	// Обработка результата
	l.procesData(order, accrual, err)

}

func (l *GatewaySrv) SendOrdersToDistantSrv(orders []*mod.Order) []*mod.Order {
	if 0 == len(orders) {
		return orders
	}

	// Массово меняем статус на "PROCESSING"
	err := l.Repo.UpdateSetStatuses2(context.TODO(), orders)
	if err != nil {
		l.Core.Logg.RaiseInfo("GatewaySrv>SendOrderToService: query is bad")
	}

	for _, order := range orders {
		url := fmt.Sprintf("%s%s%s", l.host, l.path, order.Code)

		// Ограничитель одновременно выполняемых запросов
		l.semaphore <- struct{}{}

		// Отправляем запрос в обработку
		go l.sendOrderToDistantSrv(order, url)
	}
	return orders[:0]
}
