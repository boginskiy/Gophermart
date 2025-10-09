package repository

import (
	"errors"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/models"
)

var Store map[string]*models.User

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
	user := record.(*models.User)
	Store[user.Login] = user
	return nil
}

func (r *Repo) SelectRecord(item any) (record any, err error) {
	login := item.(string)
	user, ok := Store[login]
	if !ok {
		return nil, errors.New("user bad")
	}
	return user, nil
}

// 	DeleteRecord(any) error
// 	UpdateRecord(any) error
