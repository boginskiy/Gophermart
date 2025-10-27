package service

import (
	"context"
	"encoding/json"
	"net/http"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/auth"
	"github.com/boginskiy/Gophermart/internal/logg"
	repo "github.com/boginskiy/Gophermart/internal/repository"
)

type WithdrawalsSrv struct {
	Config conf.Config
	Logger logg.Logger
	Repo   repo.RepoLoyaltyOrdersTber
}

func NewWithdrawalsSrv(config conf.Config, logger logg.Logger, repoRepoLoyaltyOrders repo.RepoLoyaltyOrdersTber) *WithdrawalsSrv {
	return &WithdrawalsSrv{
		Config: config,
		Logger: logger,
		Repo:   repoRepoLoyaltyOrders,
	}
}

func (w *WithdrawalsSrv) GetWithdrawalHistory(req *http.Request) ([]byte, error) {
	userID := req.Context().Value(auth.CtxUserID)
	withdrawals, err := w.Repo.ReadDeductions(context.TODO(), userID.(int64))

	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return EmptySliceOfBytes, nil
	}

	return json.Marshal(withdrawals)
}
