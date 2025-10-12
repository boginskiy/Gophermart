package repository

import (
	"context"
	"database/sql"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
)

type RepoUsers struct {
	Args  config.Argser
	Logg  logg.Logger
	store store.Dber
}

func NewRepoUsers(argser config.Argser, logger logg.Logger, dber store.Dber) *RepoUsers {
	return &RepoUsers{
		Args:  argser,
		Logg:  logger,
		store: dber,
	}
}

func (ru *RepoUsers) CheckUnicRecord(ctx context.Context, item any) (bool, error) {
	login, ok := item.(string)
	db, ok := ru.store.GetDB().(*sql.DB)
	if !ok {
		return false, ErrType
	}

	row := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE login = $1);`,
		login)

	var exists bool
	row.Scan(&exists)
	return exists, nil
}

func (ru *RepoUsers) InsertRecord(ctx context.Context, record any) (int64, error) {
	newUser, ok := record.(models.User)
	db, ok := ru.store.GetDB().(*sql.DB)
	if !ok {
		return 0, ErrType
	}

	row, _ := db.ExecContext(ctx,
		`INSERT INTO urls (login, password, created_at, updated_at, lastlogin_at, is_active, role)
		 VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		newUser.Login,
		newUser.Password,
		newUser.CreatedAt,
		newUser.UpdatedAt,
		newUser.LastLoginAt,
		newUser.IsActive,
		newUser.Role)

	return row.LastInsertId()

}

func (ru *RepoUsers) SelectRecord(ctx context.Context, item any) (any, error) {
	login, ok := item.(string)
	db, ok := ru.store.GetDB().(*sql.DB)
	if !ok {
		return nil, ErrType
	}

	row := db.QueryRowContext(ctx,
		`SELECT (id, login, password, created_at, updated_at, lastlogin_at, is_active, role)
		 FROM users WHERE login = $1);`,
		login)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.IsActive,
		&user.Role)
	return user, err
}

func (ru *RepoUsers) DeleteRecord(ctx context.Context, record any) error {
	return nil
}

func (ru *RepoUsers) UpdateRecord(ctx context.Context, record any) error {
	return nil
}
