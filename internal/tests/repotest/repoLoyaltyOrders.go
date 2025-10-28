package repotest

import (
	"context"

	"github.com/boginskiy/Gophermart/internal/store"
	mod "github.com/boginskiy/Gophermart/models"
)

type TestRepoLoyaltyOrders struct {
	Store store.DataBase
}

func NewTestRepoLoyaltyOrders(dataBase store.DataBase) *TestRepoLoyaltyOrders {
	return &TestRepoLoyaltyOrders{Store: dataBase}
}

func (rl *TestRepoLoyaltyOrders) CheckUnic(ctx context.Context, item any) (bool, error) {
	return false, nil
}

func (rl *TestRepoLoyaltyOrders) Create(ctx context.Context, record *mod.LoyaltyOrder) (id int64, err error) {
	db := rl.Store.GetDB().(map[string]any)
	nextNum := int64(len(db) + 1)
	record.ID = nextNum
	db[record.Code] = record
	return nextNum, nil
}

func (rl *TestRepoLoyaltyOrders) Read(ctx context.Context, item any) (record *mod.LoyaltyOrder, err error) {
	return nil, nil
}

func (rl *TestRepoLoyaltyOrders) Update(ctx context.Context, record *mod.LoyaltyOrder) error {
	return nil
}

func (rl *TestRepoLoyaltyOrders) Delete(ctx context.Context, record *mod.LoyaltyOrder) error {
	return nil
}

func (rl *TestRepoLoyaltyOrders) TotalSumOfDeductions(ctx context.Context, userID int64) (float64, error) {
	db := rl.Store.GetDB().(map[string]any)
	var totalSum float64

	for _, v := range db {
		if loyaltyOrder, ok := v.(*mod.LoyaltyOrder); ok {
			if loyaltyOrder.UserID == userID {
				totalSum += loyaltyOrder.Deduction
			}
		}
	}
	return totalSum, nil
}

func (rl *TestRepoLoyaltyOrders) ReadDeductions(ctx context.Context, userID int64) (records []*mod.LoyaltyOrder, err error) {
	db := rl.Store.GetDB().(map[string]any)
	records = make([]*mod.LoyaltyOrder, 0, 10)

	for _, v := range db {
		if loyaltyOrder, ok := v.(*mod.LoyaltyOrder); ok {
			if loyaltyOrder.UserID == userID {
				records = append(records, loyaltyOrder)
			}
		}
	}
	return records, nil
}
