package repository

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"strings"

// 	mod "github.com/boginskiy/Gophermart/models"
// )

// func UpdateSetStatuses2[T ~*mod.Order | ~*mod.Accrual](ctx context.Context, db *sql.DB, records []T) error {

// 	statuses := make([]string, 0, 30)
// 	accruals := make([]string, 0, 30)
// 	orders := make([]string, 0, 30)
// 	args := make([]any, 0, 30)

// 	// Массовая запись в БД
// 	for _, record := range records {
// 		idxCode := len(args) + 1    // Индекс кода заказа
// 		idxStatus := idxCode + 1    // Индекс статуса
// 		idxAccrual := idxStatus + 1 // Индекс суммы

// 		// Аргументы
// 		statuses = append(statuses, fmt.Sprintf("WHEN $%d THEN $%d", idxCode, idxStatus))
// 		accruals = append(accruals, fmt.Sprintf("WHEN $%d THEN CAST($%d AS INTEGER)", idxCode, idxAccrual))
// 		orders = append(orders, fmt.Sprintf("$%d", idxCode))

// 		// Параметры
// 		args = append(args, record.Order, record.Status, record.Accrual)
// 	}

// 	query := fmt.Sprintf(`UPDATE orders
// 			SET status = CASE code
// 						%s
// 					END,
// 				accrual = CASE code
// 						%s
// 					END
// 			WHERE code IN (%s)`,

// 		strings.Join(statuses, "\n"),
// 		strings.Join(accruals, "\n"),
// 		strings.Join(orders, ", "))

// 	_, err := db.ExecContext(context.TODO(), query, args...)
// 	if err != nil {
// 		rb.Logg.RaiseError("RepoOrders>UpdateSetStatuses>ExecContext", err)
// 	}

// 	return err
// }
