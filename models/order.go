package models

import "time"

type Order struct {
	ID         int64     `json:"id"`                // Primary kye
	Code       string    `json:"code"`              // Уникальный номер заказа
	Status     string    `json:"status"`            // Status заказа
	Accrual    int       `json:"accrual,omitempty"` // Положенные бонусы
	UploadedAt time.Time `json:"uploaded_at"`       // Дата время загрузки заказа
	UserID     int64     `json:"-"`                 // Связь с заказчиком
}

func NewOrder(code string, userID int64) *Order {
	return &Order{
		Code:       code,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: time.Now(),
		UserID:     userID,
	}
}

type UserOrder struct {
	Order
	Login string
}

type Accrual struct {
	Order   string `json:"order"`   // Уникальный номер заказа
	Status  string `json:"status"`  // Статус обработки
	Accrual int    `json:"accrual"` // начисленные бонусы

}
