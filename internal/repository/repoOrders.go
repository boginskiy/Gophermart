package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
	mod "github.com/boginskiy/Gophermart/models"
)

type RepoOrders struct {
	Config conf.Config
	Logger logg.Logger
	Store  store.DataBase
}

func NewRepoOrders(config conf.Config, logger logg.Logger, dataBase store.DataBase) RepoOrdersTber {
	return &RepoOrders{Config: config, Logger: logger, Store: dataBase}
}

func (rb *RepoOrders) CheckUnic(ctx context.Context, item any) (bool, error) {
	return false, nil
}

func (rb *RepoOrders) Create(ctx context.Context, record *mod.Order) (id int64, err error) {
	db := rb.Store.GetDB().(*sql.DB)

	row := db.QueryRowContext(ctx,
		`INSERT INTO orders (code, status, accrual, uploaded_at, user_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id;`,
		record.Code,
		record.Status,
		record.Accrual,
		record.UploadedAt,
		record.UserID)

	return id, row.Scan(&id)
}

func (rb *RepoOrders) Read(ctx context.Context, item any) (record *mod.Order, err error) {
	db := rb.Store.GetDB().(*sql.DB)
	orderCode, ok := item.(string)
	if !ok {
		return nil, ErrType
	}

	row := db.QueryRowContext(ctx,
		`SELECT id, code, status, accrual, uploaded_at, user_id
		 FROM orders 
		 WHERE code = $1`,
		orderCode)

	var order mod.Order

	err = row.Scan(
		&order.ID,
		&order.Code,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
		&order.UserID)

	return &order, err
}

func (rb *RepoOrders) Update(ctx context.Context, record *mod.Order) error {
	return nil
}

func (rb *RepoOrders) Delete(ctx context.Context, record *mod.Order) error {
	return nil
}

func (rb *RepoOrders) ReadWithUser(ctx context.Context, orderCode string) (record *mod.UserOrder, err error) {
	db := rb.Store.GetDB().(*sql.DB)

	row := db.QueryRowContext(ctx,
		`SELECT 

			o.id, o.code, o.status, o.accrual, o.uploaded_at, o.user_id,
		 	u.login

		 FROM orders AS o
		 INNER JOIN users AS u 
		 ON o.user_id = u.id
		 WHERE o.code = $1`,
		orderCode)

	var userOrder mod.UserOrder
	err = row.Scan(
		&userOrder.ID,
		&userOrder.Code,
		&userOrder.Status,
		&userOrder.Accrual,
		&userOrder.UploadedAt,
		&userOrder.UserID,
		&userOrder.Login,
	)
	return &userOrder, err
}

func (rb *RepoOrders) UpdateSetStatuses(ctx context.Context, records []*mod.Accrual) error {
	db := rb.Store.GetDB().(*sql.DB)

	statuses := make([]string, 0, 30)
	accruals := make([]string, 0, 30)
	orders := make([]string, 0, 30)
	args := make([]any, 0, 30)

	// Массовая запись в БД
	for _, record := range records {
		idxCode := len(args) + 1    // Индекс кода заказа
		idxStatus := idxCode + 1    // Индекс статуса
		idxAccrual := idxStatus + 1 // Индекс суммы

		// Аргументы
		statuses = append(statuses, fmt.Sprintf("WHEN $%d THEN $%d", idxCode, idxStatus))
		accruals = append(accruals, fmt.Sprintf("WHEN $%d THEN CAST($%d AS FLOAT)", idxCode, idxAccrual))
		orders = append(orders, fmt.Sprintf("$%d", idxCode))

		// Параметры
		args = append(args, record.Order, record.Status, record.Accrual)
	}

	query := fmt.Sprintf(`UPDATE orders
			SET status = CASE code
						%s
					END,
				accrual = CASE code
						%s
					END
			WHERE code IN (%s)`,

		strings.Join(statuses, "\n"),
		strings.Join(accruals, "\n"),
		strings.Join(orders, ", "))

	_, err := db.ExecContext(context.TODO(), query, args...)
	if err != nil {
		rb.Logger.RaiseError("RepoOrders>UpdateSetStatuses>ExecContext", err)
	}

	return err
}

func (rb *RepoOrders) UpdateSetStatuses2(ctx context.Context, records []*mod.Order) error {
	db := rb.Store.GetDB().(*sql.DB)

	statuses := make([]string, 0, 30)
	accruals := make([]string, 0, 30)
	orders := make([]string, 0, 30)
	args := make([]any, 0, 30)

	// Массовая запись в БД
	for _, record := range records {
		idxCode := len(args) + 1  // Индекс кода заказа
		idxAccrual := idxCode + 1 // Индекс суммы

		// Аргументы
		statuses = append(statuses, fmt.Sprintf("WHEN $%d THEN 'PROCESSING'", idxCode))
		accruals = append(accruals, fmt.Sprintf("WHEN $%d THEN CAST($%d AS FLOAT)", idxCode, idxAccrual))
		orders = append(orders, fmt.Sprintf("$%d", idxCode))

		// Параметры
		args = append(args, record.Code, record.Accrual)
	}

	query := fmt.Sprintf(`UPDATE orders
			SET status = CASE code
						%s
					END,
				accrual = CASE code
						%s
					END
			WHERE code IN (%s)`,

		strings.Join(statuses, "\n"),
		strings.Join(accruals, "\n"),
		strings.Join(orders, ", "))

	_, err := db.ExecContext(context.TODO(), query, args...)
	if err != nil {
		rb.Logger.RaiseError("RepoOrders>UpdateSetStatuses>ExecContext", err)
	}

	return err
}

func (rb *RepoOrders) ReadOrdersWithSort(ctx context.Context, userID int64) (records []*mod.Order, err error) {
	db := rb.Store.GetDB().(*sql.DB)
	rows, err := db.QueryContext(ctx,
		`SELECT
			id, code, status, accrual, uploaded_at, user_id
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC`,
		userID)

	if err != nil || rows.Err() != nil {
		rb.Logger.RaiseError("RepoOrders>ReadOrdersWithSort>QueryContext", rows.Err())
		return nil, err
	}

	defer rows.Close()
	records = make([]*mod.Order, 0, 10)

	for rows.Next() {
		var record mod.Order
		err := rows.Scan(
			&record.ID,
			&record.Code,
			&record.Status,
			&record.Accrual,
			&record.UploadedAt,
			&record.UserID)

		if err != nil {
			rb.Logger.RaiseInfo(err.Error())
		} else {
			records = append(records, &record)
		}
	}
	return records, nil
}

func (rb *RepoOrders) ReadAccruals(ctx context.Context, userID int64) (float64, error) {
	db := rb.Store.GetDB().(*sql.DB)
	var totalSum float64

	err := db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(accrual), 0.0)::float
		 FROM orders
		 WHERE user_id = $1`,
		userID).Scan(&totalSum)
	if err != nil {
		return 0, err
	}
	return totalSum, nil
}
