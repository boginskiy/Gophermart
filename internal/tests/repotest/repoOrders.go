package repotest

import (
	"context"
	"errors"

	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/internal/store"
	"github.com/boginskiy/Gophermart/models"
)

type TestRepoOrders struct {
	Store store.DataBase
}

func NewTestRepoOrders(dataBase store.DataBase) *TestRepoOrders {
	return &TestRepoOrders{Store: dataBase}
}

func (ru *TestRepoOrders) CheckUnic(ctx context.Context, item any) (bool, error) {
	return false, nil
}

func (ru *TestRepoOrders) Create(ctx context.Context, record *models.Order) (id int64, err error) {
	db := ru.Store.GetDB().(map[string]any)
	nextNum := int64(len(db) + 1)
	record.ID = nextNum
	db[record.Code] = record
	return nextNum, nil
}

func (ru *TestRepoOrders) Read(ctx context.Context, item any) (record *models.Order, err error) {
	db := ru.Store.GetDB().(map[string]any)
	code, ok := item.(string)
	if !ok {
		return nil, repository.ErrType
	}

	order, ok := db[code]
	if !ok {
		return nil, errors.New("empty")
	}
	return order.(*models.Order), nil
}

func (ru *TestRepoOrders) Update(ctx context.Context, record *models.Order) error {
	return nil
}

func (ru *TestRepoOrders) Delete(ctx context.Context, record *models.Order) error {
	return nil
}

func (ru *TestRepoOrders) ReadWithUser(ctx context.Context, code string) (*models.UserOrder, error) {
	db := ru.Store.GetDB().(map[string]any)
	var userOrder models.UserOrder

	record, ok := db[code]
	if !ok {
		return &userOrder, errors.New("data is empty")
	}

	order := record.(*models.Order)

	// Собираем *models.UserOrder
	userOrder.ID = order.ID
	userOrder.Code = order.Code
	userOrder.Status = order.Status
	userOrder.Accrual = order.Accrual
	userOrder.UploadedAt = order.UploadedAt
	userOrder.UserID = order.UserID

	// Ищем пользователя
	for _, v := range db {
		if user, ok := v.(*models.User); ok {
			if user.ID == order.UserID {
				userOrder.Login = user.Login
				return &userOrder, nil
			}
		}
	}
	return nil, errors.New("data is bad")
}

func (ru *TestRepoOrders) UpdateSetStatuses(ctx context.Context, records []*models.Accrual) error {
	return nil
}

func (ru *TestRepoOrders) UpdateSetStatuses2(ctx context.Context, records []*models.Order) error {
	return nil
}

func (ru *TestRepoOrders) ReadOrdersWithSort(ctx context.Context, userID int64) (records []*models.Order, err error) {
	db := ru.Store.GetDB().(map[string]any)
	records = make([]*models.Order, 0, 10)

	for _, record := range db {
		if order, ok := record.(*models.Order); ok {
			if order.UserID == userID {
				records = append(records, order)
			}
		}
	}
	return records, nil
}

func (ru *TestRepoOrders) ReadAccruals(ctx context.Context, userID int64) (float64, error) {
	db := ru.Store.GetDB().(map[string]any)
	var totalSum float64

	for _, v := range db {
		if order, ok := v.(*models.Order); ok {
			if order.UserID == userID {
				totalSum += order.Accrual
			}
		}
	}
	return totalSum, nil
}
