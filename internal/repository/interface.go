package repository

import (
	"context"

	mod "github.com/boginskiy/Gophermart/models"
)

type RepoCRUDer[T any] interface {
	CheckUnic(ctx context.Context, item any) (bool, error)
	Create(ctx context.Context, model *T) (id int64, err error)
	Read(ctx context.Context, item any) (model *T, err error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, model *T) error
}

type RepoOrdersTber interface {
	RepoCRUDer[mod.Order]
	// Расширение CRUD интерфейса
	ReadWithUser(ctx context.Context, orderCode string) (record *mod.UserOrder, err error)
	UpdateSetStatuses(ctx context.Context, records []*mod.Accrual) error
	UpdateSetStatuses2(ctx context.Context, records []*mod.Order) error
	ReadOrdersWithSort(ctx context.Context, userID int64) (records []*mod.Order, err error)
	ReadAccruals(ctx context.Context, userID int64) (int, error)
}

type RepoUsersTber interface {
	RepoCRUDer[mod.User]
	// Расширение CRUD интерфейса
}

type RepoLoyaltyOrdersTber interface {
	RepoCRUDer[mod.LoyaltyOrder]
	// Расширение CRUD интерфейса
	ReadDeductions(ctx context.Context, userID int64) (int, error)
}
