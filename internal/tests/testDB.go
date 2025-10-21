package tests

import (
	"log"

	"github.com/boginskiy/Gophermart/models"
)

type TestDB struct {
	db map[string]any
}

func NewTestDB() *TestDB {

	login := "TestMan" // Тестовый user
	user, err := models.NewUser(login, "123456")
	user.ID = int64(1)
	if err != nil {
		log.Fatal(err)
	}

	code := "68379148" // Тестовый order
	order := models.NewOrder(code, int64(1))
	order.UserID = user.ID

	store := map[string]any{login: user, code: order}
	return &TestDB{db: store}
}

func (tDB *TestDB) Open() {
}

func (tDB *TestDB) Clouse() {
}

func (tDB *TestDB) Ping() {
}

func (tDB *TestDB) GetDB() any {
	return tDB.db
}
