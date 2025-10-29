package service

import (
	"context"
	"encoding/json"
	"net/http"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/logg"
	repo "github.com/boginskiy/Gophermart/internal/repository"
	mod "github.com/boginskiy/Gophermart/models"
)

type BalanceServ struct {
	Config            conf.Config
	Logger            logg.Logger
	RepoOrders        repo.RepoOrdersTber
	RepoLoyaltyOrders repo.RepoLoyaltyOrdersTber
	OrderCheck        OrderChecker
}

func NewBalanceServ(
	config conf.Config,
	logger logg.Logger,
	repoOrders repo.RepoOrdersTber,
	repoLoyaltyOrders repo.RepoLoyaltyOrdersTber,
	orderCheck OrderChecker) *BalanceServ {

	return &BalanceServ{
		Config:            config,
		Logger:            logger,
		RepoOrders:        repoOrders,
		RepoLoyaltyOrders: repoLoyaltyOrders,
		OrderCheck:        orderCheck,
	}
}

func (bs *BalanceServ) GetBalance(req *http.Request) ([]byte, error) {
	userID := req.Context().Value(auth.CtxUserID)

	accruals, err := bs.RepoOrders.ReadAccruals(context.TODO(), userID.(int64))
	if err != nil {
		return nil, err
	}

	deductions, err := bs.RepoLoyaltyOrders.TotalSumOfDeductions(context.TODO(), userID.(int64))
	if err != nil {
		return nil, err
	}

	var statusBalance mod.StatusBalance
	statusBalance.Current = accruals - deductions
	statusBalance.Withdrawn = deductions

	return json.Marshal(statusBalance)
}

func (bs *BalanceServ) GetWithdrawal(req *http.Request) ([]byte, error) {
	userID := req.Context().Value(auth.CtxUserID)

	// Парсим данные request
	loyaltyOrder := mod.NewLoyaltyOrder("", 0, userID.(int64))
	json.NewDecoder(req.Body).Decode(loyaltyOrder)

	// Проверка номера заказа алгоритмом Luna || номер заказа не должен быть пустым
	if !bs.OrderCheck.CheckDigits(loyaltyOrder.Code) || loyaltyOrder.Code == "" {
		return nil, ErrOrderNumber
	}

	// Заглядываем в БД и смотрим сколько бонусов накопилось
	accruals, err := bs.RepoOrders.ReadAccruals(context.TODO(), userID.(int64))
	if err != nil {
		return nil, err
	}

	// Делаем расчет остатка после запроса на бонус и если он отрицательный
	// то пользователь запросил больше, чем может унести
	bonuse := accruals - loyaltyOrder.Deduction
	if bonuse < 0 {
		return nil, ErrBonuseLimit
	}

	// Фиксируем списание бонусов в БД
	_, err = bs.RepoLoyaltyOrders.Create(context.TODO(), loyaltyOrder)
	if err != nil {
		return nil, err
	}

	return json.Marshal(loyaltyOrder)
}
