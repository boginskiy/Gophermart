package repository

import (
	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type Repo struct {
	Args config.Argser
	Logg logg.Logger
}

func NewRepos(argser config.Argser, logger logg.Logger) *Repo {
	return &Repo{
		Args: argser,
		Logg: logger,
	}
}

func (r *Repo) CheckUnicRecord(record any) bool {
	return true
}

func (r *Repo) InsertRecord(record any) error {
	return nil
}

// CheckUnicRecord(any) error
// 	InsertRecord(any) error
// 	SelectRecord(any) error
// 	DeleteRecord(any) error
// 	UpdateRecord(any) error
