package models

import "time"

type Order struct {
	ID         int       `json:"id"`                // Уникальный номер заказа
	Status     string    `json:"status"`            // Имя пользователя (может содержать ФИО)
	Accrual    int       `json:"accrual,omitempty"` // Положенные бонусы
	UploadedAt time.Time `json:"uploaded_at"`       // Дата время загрузки заказа
	UserID     int       `json:"-"`                 // Связь с заказчиком
}
