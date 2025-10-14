package repository

import (
	"context"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
)

type RepoBalances struct {
	Args  config.Argser
	Logg  logg.Logger
	Store store.Dber
}

func NewRepoBalances(argser config.Argser, logger logg.Logger, dber store.Dber) RepoBalancesTber {
	return &RepoBalances{Args: argser, Logg: logger, Store: dber}
}

func (rb *RepoBalances) CheckUnic(ctx context.Context, item any) (bool, error) {
	return false, nil
}

func (rb *RepoBalances) Create(ctx context.Context, record *models.Balance) (id int64, err error) {
	return 0, nil
}

func (rb *RepoBalances) Read(ctx context.Context, item any) (record *models.Balance, err error) {
	return nil, nil
}

func (rb *RepoBalances) Update(ctx context.Context, record *models.Balance) error {
	return nil
}

func (rb *RepoBalances) Delete(ctx context.Context, record *models.Balance) error {
	return nil
}
