package repository

import (
	"context"
	"database/sql"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
)

type RepoUsers struct {
	Config conf.Config
	Logger logg.Logger
	Store  store.DataBase
}

func NewRepoUsers(config conf.Config, logger logg.Logger, dataBase store.DataBase) RepoUsersTber {
	return &RepoUsers{Config: config, Logger: logger, Store: dataBase}
}

func (ru *RepoUsers) CheckUnic(ctx context.Context, item any) (bool, error) {
	db := ru.Store.GetDB().(*sql.DB)
	login, ok := item.(string)
	if !ok {
		return false, ErrType
	}

	row := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE login = $1);`,
		login)

	var exists bool
	row.Scan(&exists)
	return !exists, nil
}

func (ru *RepoUsers) Create(ctx context.Context, record *models.User) (id int64, err error) {
	db := ru.Store.GetDB().(*sql.DB)

	row := db.QueryRowContext(ctx,
		`INSERT INTO users (login, password, created_at, updated_at, lastlogin_at, is_active, role)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id;`,
		record.Login,
		record.Password,
		record.CreatedAt,
		record.UpdatedAt,
		record.LastLoginAt,
		record.IsActive,
		record.Role)

	return id, row.Scan(&id)
}

func (ru *RepoUsers) Read(ctx context.Context, item any) (record *models.User, err error) {
	db := ru.Store.GetDB().(*sql.DB)
	login, ok := item.(string)
	if !ok {
		return nil, ErrType
	}

	row := db.QueryRowContext(ctx,
		`SELECT id, login, password, created_at, updated_at, lastlogin_at, is_active, role
		 FROM users 
		 WHERE login = $1`,
		login)

	var user models.User

	err = row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.IsActive,
		&user.Role)

	return &user, err
}

func (ru *RepoUsers) Update(ctx context.Context, record *models.User) error {
	return nil
}

func (ru *RepoUsers) Delete(ctx context.Context, record *models.User) error {
	return nil
}
