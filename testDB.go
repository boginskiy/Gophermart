package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/boginskiy/Gophermart/models"
)

func OpenDB(args string) (*sql.DB, error) {
	db, err := sql.Open("postgres", args)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	// DB
	db, err := OpenDB("postgres://username:userpassword@localhost:5432/gophermartdb?sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	// User
	user, _ := models.NewUser("Vasiy2", "1234")
	//
	// newOrder := models.NewOrder("12345678", 1)

	row := db.QueryRowContext(context.TODO(),
		`INSERT INTO users (login, password, created_at, updated_at, lastlogin_at, is_active, role)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id;`,
		user.Login,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLoginAt,
		user.IsActive,
		user.Role)

	var userID int64
	row.Scan(&userID)

	fmt.Println(userID)

	return
}
