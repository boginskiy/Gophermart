package repository

import "context"

type RepoDBer interface {
	CheckUnic(repoTb RepoTber, item any) bool
	Create(repoTb RepoTber, record any) error
	Read(repoTb RepoTber, item any) (record any, err error)
	Update(repoTb RepoTber, record any) error
	Delete(repoTb RepoTber, record any) error
}

type RepoTber interface {
	CheckUnicRecord(ctx context.Context, item any) (bool, error)
	InsertRecord(ctx context.Context, record any) (int64, error)
	SelectRecord(ctx context.Context, item any) (record any, err error)
	DeleteRecord(ctx context.Context, record any) error
	UpdateRecord(ctx context.Context, record any) error
}
