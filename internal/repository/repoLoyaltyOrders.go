package repository

import (
	"context"
	"database/sql"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
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

func (rl *RepoLoyaltyOrders) Create(ctx context.Context, record *models.LoyaltyOrder) (id int64, err error) {
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

func (rl *RepoLoyaltyOrders) Read(ctx context.Context, item any) (record *models.LoyaltyOrder, err error) {
	return nil, nil
}

func (rl *RepoLoyaltyOrders) Update(ctx context.Context, record *models.LoyaltyOrder) error {
	return nil
}

func (rl *RepoLoyaltyOrders) Delete(ctx context.Context, record *models.LoyaltyOrder) error {
	return nil
}

func (rl *RepoLoyaltyOrders) ReadDeductions(ctx context.Context, userID int64) (int, error) {
	db := rl.Store.GetDB().(*sql.DB)
	var totalSum int

	err := db.QueryRowContext(ctx,
		`SELECT SUM(deduction)
		 FROM loyalty_orders
		 WHERE user_id = $1`,
		userID).Scan(&totalSum)
	if err != nil {
		return 0, err
	}
	return totalSum, nil
}
