package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	repo "github.com/boginskiy/Gophermart/internal/repository"
	mod "github.com/boginskiy/Gophermart/models"
)

type BalanceServ struct {
	Core              *CoreSrv
	RepoOrders        repo.RepoOrdersTber
	RepoLoyaltyOrders repo.RepoLoyaltyOrdersTber
}

func NewBalanceServ(c *CoreSrv, repoOrders repo.RepoOrdersTber, repoLoyaltyOrders repo.RepoLoyaltyOrdersTber) *BalanceServ {
	return &BalanceServ{
		Core:              c,
		RepoOrders:        repoOrders,
		RepoLoyaltyOrders: repoLoyaltyOrders,
	}
}

func (bs *BalanceServ) GetBalance(req *http.Request) ([]byte, error) {
	userID := bs.Core.takeParamFromAuth(req, auth.CtxUserID)

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
	userID := bs.Core.takeParamFromAuth(req, auth.CtxUserID)

	// Парсим данные request
	loyaltyOrder := mod.NewLoyaltyOrder("", 0, userID.(int64))
	err := json.NewDecoder(req.Body).Decode(loyaltyOrder)

	if err != nil {
		return nil, err
	}

	// Проверка номера заказа алгоритмом Luna
	if !bs.Core.OrderCheck.CheckDigits(loyaltyOrder.Code) {
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
