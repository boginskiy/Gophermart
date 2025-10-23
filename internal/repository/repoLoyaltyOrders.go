package repository

import (
	"context"
	"database/sql"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	mod "github.com/boginskiy/Gophermart/models"
)

type RepoLoyaltyOrders struct {
	Args  config.Argser
	Logg  logg.Logger
	Store store.Dber
}

func NewRepoLoyaltyOrders(argser config.Argser, logger logg.Logger, dber store.Dber) RepoLoyaltyOrdersTber {
	return &RepoLoyaltyOrders{Args: argser, Logg: logger, Store: dber}
}

func (rl *RepoLoyaltyOrders) CheckUnic(ctx context.Context, item any) (bool, error) {
	return false, nil
}

func (rl *RepoLoyaltyOrders) Create(ctx context.Context, record *mod.LoyaltyOrder) (id int64, err error) {
	db := rl.Store.GetDB().(*sql.DB)

	row := db.QueryRowContext(ctx,
		`INSERT INTO loyalty_orders (code, deduction, processed_at, user_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id;`,
		record.Code,
		record.Deduction,
		record.ProcessedAt,
		record.UserID)

	return id, row.Scan(&id)
}

func (rl *RepoLoyaltyOrders) Read(ctx context.Context, item any) (record *mod.LoyaltyOrder, err error) {
	return nil, nil
}

func (rl *RepoLoyaltyOrders) Update(ctx context.Context, record *mod.LoyaltyOrder) error {
	return nil
}

func (rl *RepoLoyaltyOrders) Delete(ctx context.Context, record *mod.LoyaltyOrder) error {
	return nil
}

func (rl *RepoLoyaltyOrders) TotalSumOfDeductions(ctx context.Context, userID int64) (float64, error) {
	db := rl.Store.GetDB().(*sql.DB)
	var totalSum float64

	err := db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(deduction), 0.0)::float
		 FROM loyalty_orders
		 WHERE user_id = $1`,
		userID).Scan(&totalSum)

	if err != nil {
		return 0, err
	}
	return totalSum, nil
}

func (rl *RepoLoyaltyOrders) ReadDeductions(ctx context.Context, userID int64) (records []*mod.LoyaltyOrder, err error) {
	db := rl.Store.GetDB().(*sql.DB)
	rows, err := db.QueryContext(ctx,
		`SELECT
			id, code, deduction, processed_at, user_id
		FROM loyalty_orders
		WHERE user_id = $1
		ORDER BY processed_at DESC`,
		userID)

	if err != nil || rows.Err() != nil {
		rl.Logg.RaiseError("RepoLoyaltyOrders>ReadDeductions>QueryContext", rows.Err())
		return nil, err
	}

	defer rows.Close()
	records = make([]*mod.LoyaltyOrder, 0, 10)

	for rows.Next() {
		var record mod.LoyaltyOrder
		err := rows.Scan(
			&record.ID,
			&record.Code,
			&record.Deduction,
			&record.ProcessedAt,
			&record.UserID)

		if err != nil {
			rl.Logg.RaiseInfo(err.Error())
		} else {
			records = append(records, &record)
		}
	}
	return records, nil
}
