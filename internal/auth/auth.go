package auth

import (
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/models"
)

type Auth struct {
	Repo    repository.Repository
	JWTServ JWTokener
	Core    *AhCore
}

func NewAuth(core *AhCore, jwter JWTokener, repo repository.Repository) *Auth {
	return &Auth{
		Repo:    repo,
		Core:    core,
		JWTServ: jwter,
	}
}

func (u *Auth) Authentication(req *http.Request) ([]byte, error) {
	return nil, nil
}

func (u *Auth) Registration(req *http.Request) ([]byte, *http.Cookie, error) {
	// Достаем login, password
	login, password, err := u.Core.takeLoginAndPassword(req)
	if err != nil {
		return nil, nil, err
	}

	// Проверка уникальности login
	loginIsUnic := u.Repo.CheckUnicRecord(login)
	if !loginIsUnic {
		// Пользователь не уникален, логин уже занят
		return nil, nil, ErrLoginNotUnic
	}

	// Создаем нового пользователя
	newUser, err := models.NewUser(login, password)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>NewUser", err)
		return nil, nil, ErrCreateUser
	}

	// Запись нового пользователя в БД
	err = u.Repo.InsertRecord(newUser)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>InsertRecord", err)
		return nil, nil, ErrCreateUser
	}

	// Выдаем токен для последующей аутентификации
	// Выдаем токен по логину
	// dd, r := req.Cookie()

	// Ответ
	userByte, err := json.Marshal(newUser)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>Marshal", err)
		return nil, nil, ErrCreateUser
	}

	return userByte, nil, nil
}
