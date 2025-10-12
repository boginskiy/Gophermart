package service

import (
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/repository"
)

type UserSrv struct {
	Repo   repository.RepoDBer
	Logger logg.Logger
}

func NewUserSrv(repo repository.RepoDBer, logger logg.Logger) *UserSrv {
	return &UserSrv{
		Repo:   repo,
		Logger: logger,
	}
}
