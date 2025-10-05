package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/models"
)

type UserSrv struct {
}

func NewUserSrv() *UserSrv {
	return &UserSrv{}
}

func (u *UserSrv) getLoginAndPassword(st map[string]any) (login, password string) {
	if login, ok := st["login"].(string); ok {
		if password, ok2 := st["password"].(string); ok2 {
			return login, password
		}
	}
	return "", ""
}

func (u *UserSrv) Registration(req *http.Request) ([]byte, error) {
	// Читаем тело запроса
	dataByte, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// Парсинг данных
	overallStruct := make(map[string]any, 2)
	json.Unmarshal(dataByte, &overallStruct)

	login, password := u.getLoginAndPassword(overallStruct)
	if password == "" {
		return nil, errors.New("user not registration")
	}

	// Создаем нового пользователя
	user, err := models.NewUser(login, password)
	if err != nil {
		return nil, errors.New("user not create")
	}

	// Записываем в БД TODO!
	u.Repo.InsertRecord(user)

	userByte, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New("user not serialize")
	}

	return userByte, nil
}

func (u *UserSrv) Authentication(req *http.Request) ([]byte, error) {
	return []byte{}, errors.New("")
}
