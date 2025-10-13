package repository

import (
	"context"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
)

type RepoOrders struct {
	Args  config.Argser
	Logg  logg.Logger
	Store store.Dber
}

func NewRepoOrders(argser config.Argser, logger logg.Logger, dber store.Dber) RepoOrdersTber {
	return &RepoOrders{Args: argser, Logg: logger, Store: dber}
}

func (rb *RepoOrders) CheckUnic(ctx context.Context, item any) (bool, error) {
	return false, nil
}

func (rb *RepoOrders) Create(ctx context.Context, record *models.Order) error {
	return nil
}

func (rb *RepoOrders) Read(ctx context.Context, item any) (record *models.Order, err error) {
	return nil, nil
}

func (rb *RepoOrders) Update(ctx context.Context, record *models.Order) error {
	return nil
}

func (rb *RepoOrders) Delete(ctx context.Context, record *models.Order) error {
	return nil
}
