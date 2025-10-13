package repository

import (
	"context"

	"github.com/boginskiy/Gophermart/models"
)

type RepoCRUDer[T any] interface {
	CheckUnic(ctx context.Context, item any) (bool, error)
	Create(ctx context.Context, model *T) error
	Read(ctx context.Context, item any) (model *T, err error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, model *T) error
}

type RepoOrdersTber interface {
	RepoCRUDer[models.Order]
	// Расширение CRUD интерфейса
}

type RepoUsersTber interface {
	RepoCRUDer[models.User]
	// Расширение CRUD интерфейса
}

type RepoBalancesTber interface {
	RepoCRUDer[models.Balance]
	// Расширение CRUD интерфейса
}
