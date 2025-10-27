package repotest

import (
	"context"
	"errors"

	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
)

type TestRepoUsers struct {
	Store store.DataBase
}

func NewTestRepoUsers(dataBase store.DataBase) *TestRepoUsers {
	return &TestRepoUsers{Store: dataBase}
}

func (ru *TestRepoUsers) CheckUnic(ctx context.Context, item any) (bool, error) {
	db := ru.Store.GetDB().(map[string]any)
	login, ok := item.(string)
	if !ok {

		return false, repository.ErrType
	}
	_, ok = db[login]
	if ok {
		return false, nil
	}
	return true, nil
}

func (ru *TestRepoUsers) Create(ctx context.Context, record *models.User) (id int64, err error) {
	db := ru.Store.GetDB().(map[string]any)
	nextNum := int64(len(db) + 1)
	record.ID = nextNum
	db[record.Login] = record
	return nextNum, nil
}

func (ru *TestRepoUsers) Read(ctx context.Context, item any) (record *models.User, err error) {
	db := ru.Store.GetDB().(map[string]any)
	login, ok := item.(string)
	if !ok {
		return nil, repository.ErrType
	}

	user, ok := db[login]
	if !ok {
		return nil, errors.New("empty")
	}
	return user.(*models.User), nil
}

func (ru *TestRepoUsers) Update(ctx context.Context, record *models.User) error {
	return nil
}

func (ru *TestRepoUsers) Delete(ctx context.Context, record *models.User) error {
	return nil
}
