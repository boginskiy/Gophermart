package tests

import (
	"log"

	"github.com/boginskiy/Gophermart/models"
)

type TestDB struct {
	db map[string]any
}

func NewTestDB() *TestDB {

	login := "TestMan" // Тестовый user1
	user, err := models.NewUser(login, "123456")
	user.ID = int64(1)
	if err != nil {
		log.Fatal(err)
	}

	login2 := "TestMan2" // Тестовый user2
	user2, err := models.NewUser(login, "7890123")
	user2.ID = int64(2)
	if err != nil {
		log.Fatal(err)
	}

	code := "68379148" // Тестовый order
	order := models.NewOrder(code, int64(1))
	order.Accrual = 100.0
	order.UserID = user.ID

	keyLOrder1 := "LoyaltyOrder1" // Тестовый Loyalty Order1
	deduction1 := 25.0
	loyaltyOrder1 := models.NewLoyaltyOrder(code, deduction1, user.ID)

	keyLOrder2 := "LoyaltyOrder2" // Тестовый Loyalty Order2
	deduction2 := 5.0
	loyaltyOrder2 := models.NewLoyaltyOrder(code, deduction2, user.ID)

	store := map[string]any{
		login:      user,
		login2:     user2,
		code:       order,
		keyLOrder1: loyaltyOrder1,
		keyLOrder2: loyaltyOrder2}
	return &TestDB{db: store}
}

func (tDB *TestDB) Open() {
}

func (tDB *TestDB) Close() {
}

func (tDB *TestDB) Ping() {
}

func (tDB *TestDB) GetDB() any {
	return tDB.db
}
