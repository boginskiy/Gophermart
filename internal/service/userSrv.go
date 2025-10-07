package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/models"
)

type UserSrv struct {
	Repo   repository.Repository
	Logger logg.Logger
}

func NewUserSrv(repo repository.Repository, logger logg.Logger) *UserSrv {
	return &UserSrv{
		Repo:   repo,
		Logger: logger,
	}
}

func (u *UserSrv) takeLoginAndPassword(req *http.Request) (login, password string, err error) {
	// Читаем body req
	dataByte, err := io.ReadAll(req.Body)
	if err != nil {
		u.Logger.RaiseError("UserSrv>takeLoginAndPassword>ReadAll", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Парсинг данных
	overallStruct := make(map[string]any, 2)
	err = json.Unmarshal(dataByte, &overallStruct)
	if err != nil {
		u.Logger.RaiseError("UserSrv>takeLoginAndPassword>Unmarshal", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Достаем login, password
	if login, ok := overallStruct["login"].(string); ok {
		if password, ok2 := overallStruct["password"].(string); ok2 {
			return login, password, nil
		}
	}

	u.Logger.RaiseInfo("UserSrv>takeLoginAndPassword>data not found")
	return "", "", ErrLoginPasswordIsBad
}

func (u *UserSrv) Registration(req *http.Request) ([]byte, error) {
	// Достаем login, password
	login, password, err := u.takeLoginAndPassword(req)
	if err != nil {
		return nil, err
	}

	// Проверка уникальности login
	loginIsUnic := u.Repo.CheckUnicRecord(login)
	if !loginIsUnic {
		// Пользователь не уникален, логин уже занят
		return nil, ErrLoginNotUnic
	}

	// Создаем нового пользователя
	newUser, err := models.NewUser(login, password)
	if err != nil {
		u.Logger.RaiseError("UserSrv>Registration>NewUser", err)
		return nil, ErrCreateUser
	}

	// Запись нового пользователя в БД
	err = u.Repo.InsertRecord(newUser)
	if err != nil {
		u.Logger.RaiseError("UserSrv>Registration>InsertRecord", err)
		return nil, ErrCreateUser
	}

	// Выдаем токен для последующей аутентификации
	// Выдаем токен по логину

	// Ответ
	userByte, err := json.Marshal(newUser)
	if err != nil {
		u.Logger.RaiseError("UserSrv>Registration>Marshal", err)
		return nil, ErrCreateUser
	}

	return userByte, nil
}

func (u *UserSrv) Authentication(req *http.Request) ([]byte, error) {
	return []byte{}, errors.New("")
}
