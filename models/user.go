package models

import (
	"time"

	"github.com/boginskiy/Gophermart/pkg"
)

type User struct {
	ID          int       `json:"-"`          // Уникальный идентификатор пользователя
	Login       string    `json:"login"`      // Имя пользователя (может содержать ФИО)
	Password    string    `json:"-"`          // Хешированный пароль (не отображается при сериализации JSON)
	CreatedAt   time.Time `json:"created_at"` // Дата регистрации пользователя
	UpdatedAt   time.Time `json:"-"`          // Последнее обновление профиля
	LastLoginAt time.Time `json:"-"`          // Время последнего входа
	IsActive    bool      `json:"-"`          // Флаг активности пользователя (заблокирован или активен)
	Role        string    `json:"-"`          // Роли пользователя (администратор, модератор, обычный пользователь)
}

func NewUser(login, password string) (*User, error) {
	// Хеширование пароля
	hash, err := pkg.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Login:       login,
		Password:    hash,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		LastLoginAt: time.Now().UTC(),
		IsActive:    true,
		Role:        "user",
	}, nil
}
