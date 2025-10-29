package models

import "time"

type LoyaltyOrder struct {
	ID          int64     `json:"-"`
	Code        string    `json:"order"`         // Code is Уникальный номер заказа
	Deduction   float64   `json:"sum,omitempty"` // Deduction is Списанные бонусы в счет текущего заказа
	ProcessedAt time.Time `json:"processed_at"`  // ProcessedAt is Дата время загрузки заказа
	UserID      int64     `json:"-"`             // UserID is Связь с заказчиком
}

func NewLoyaltyOrder(code string, deduction float64, userID int64) *LoyaltyOrder {
	return &LoyaltyOrder{
		Code:        code,
		Deduction:   deduction,
		ProcessedAt: time.Now(),
		UserID:      userID}
}

type StatusBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
