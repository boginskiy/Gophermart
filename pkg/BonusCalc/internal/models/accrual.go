package models

type Accrual struct {
	Order   string `json:"order"`   // Уникальный номер заказа
	Status  string `json:"status"`  // Статус обработки
	Accrual int    `json:"accrual"` // начисленные бонусы
}
