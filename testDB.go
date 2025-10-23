package main

import "database/sql"

// import (
// 	"database/sql"

// 	_ "github.com/lib/pq"
// )

func OpenDB(args string) (*sql.DB, error) {
	db, err := sql.Open("postgres", args)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// func main() {

// 	// DB
// 	db, err := OpenDB("postgres://username:userpassword@localhost:5432/gophermartdb?sslmode=disable")
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var record1 = &models.Accrual{Order: "28561561", Status: "K", Accrual: 500}
// 	// var record2 = &models.Accrual{Order: "74701145", Status: "G", Accrual: 1000}

// 	records := []*models.Accrual{record1}

// 	statuses := make([]string, 0, 30)
// 	accruals := make([]string, 0, 30)
// 	orders := make([]string, 0, 30)
// 	args := make([]any, 0, 30)

// 	// Массовая запись в БД
// 	for _, record := range records {
// 		idxCode := len(args) + 1 // Индекс кода заказа
// 		// idxStatus := idxCode + 1    // Индекс статуса
// 		idxAccrual := idxCode + 1 // Индекс суммы

// 		// Аргументы
// 		statuses = append(statuses, fmt.Sprintf("WHEN $%d THEN 'PROCESSING'", idxCode))
// 		accruals = append(accruals, fmt.Sprintf("WHEN $%d THEN CAST($%d AS INTEGER)", idxCode, idxAccrual))
// 		orders = append(orders, fmt.Sprintf("$%d", idxCode))

// 		// Параметры
// 		args = append(args, record.Order, record.Accrual)
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

// 	fmt.Println(query)

// 	_, err = db.ExecContext(context.TODO(), query, args...)

// 	fmt.Println(err)

// }
