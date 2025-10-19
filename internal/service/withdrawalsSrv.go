package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/auth"
	repo "github.com/boginskiy/Gophermart/internal/repository"
)

type WithdrawalsSrv struct {
	Core *CoreSrv
	Repo repo.RepoLoyaltyOrdersTber
}

func NewWithdrawalsSrv(c *CoreSrv, repo repo.RepoLoyaltyOrdersTber) *WithdrawalsSrv {
	return &WithdrawalsSrv{
		Core: c,
		Repo: repo,
	}
}

func (w *WithdrawalsSrv) GetWithdrawalHistory(req *http.Request) ([]byte, error) {
	userID := w.Core.takeParamFromAuth(req, auth.CtxUserID)
	withdrawals, err := w.Repo.ReadDeductions(context.TODO(), userID.(int64))

	if err != nil {
		return nil, err
	}

	if 0 == len(withdrawals) {
		return EmptySliceOfBytes, nil
	}

	return json.Marshal(withdrawals)
}
