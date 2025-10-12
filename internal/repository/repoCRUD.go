package repository

import (
	"context"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type Repo struct {
	Args   config.Argser
	Logg   logg.Logger
	RepoTb RepoTber
}

func NewRepo(argser config.Argser, logger logg.Logger, repoTber RepoTber) *Repo {
	return &Repo{
		Args:   argser,
		Logg:   logger,
		RepoTb: repoTber,
	}
}

func (r *Repo) CheckUnic(repoTb RepoTber, item any) bool {
	result, err := repoTb.CheckUnicRecord(context.TODO(), item)
	if err != nil {
		r.Logg.RaiseError("Repo>CheckUnic>CheckUnicRecord", err)
	}
	return result
}

func (r *Repo) Create(repoTb RepoTber, record any) error {
	// TODO. Есть id вставки
	_, err := repoTb.InsertRecord(context.TODO(), record)
	if err != nil {
		r.Logg.RaiseError("Repo>Create>InsertRecord", err)
	}
	return err
}

func (r *Repo) Read(repoTb RepoTber, item any) (record any, err error) {
	record, err = repoTb.SelectRecord(context.TODO(), item)
	if err != nil {
		r.Logg.RaiseError("Repo>Read>SelectRecord", err)
	}
	return record, err
}

func (r *Repo) Update(repoTb RepoTber, record any) error {
	return repoTb.UpdateRecord(context.TODO(), record)
}

func (r *Repo) Delete(repoTb RepoTber, record any) error {
	return repoTb.DeleteRecord(context.TODO(), record)
}
