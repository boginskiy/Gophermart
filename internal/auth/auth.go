package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/models"
	"github.com/boginskiy/Gophermart/pkg"
)

type Auth struct {
	Repo    repository.RepoUsersTber
	JWTServ JWTokener
	Core    *AhCore
}

func NewAuth(core *AhCore, jwter JWTokener, repoUsers repository.RepoUsersTber) *Auth {
	return &Auth{
		Repo:    repoUsers,
		JWTServ: jwter,
		Core:    core,
	}
}

func (u *Auth) CheckToken(token string) (login, role string, err error) {
	return u.JWTServ.CheckOfValidToken(token)
}

func (u *Auth) CheckAuthReq(req *http.Request) bool {
	url := req.URL.String()
	foundR := strings.Contains(url, "register")
	if !foundR {
		foundL := strings.Contains(url, "login")
		return foundL
	}
	return foundR
}

func (u *Auth) Registration(req *http.Request) ([]byte, *http.Cookie, error) {
	// Достаем login, password
	login, password, err := u.Core.takeLoginAndPassword(req)
	if err != nil {
		return nil, nil, ErrLoginPasswordIsBad
	}

	// Проверка уникальности login
	loginIsUnic, _ := u.Repo.CheckUnic(context.TODO(), login)
	if !loginIsUnic {
		// Пользователь не уникален, логин уже занят
		u.Core.Logg.RaiseInfo(ErrLoginNotUnic.Error())
		return nil, nil, ErrLoginNotUnic
	}

	// Создаем нового пользователя
	newUser, err := models.NewUser(login, password)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>NewUser", err)
		return nil, nil, ErrCreateUser
	}

	// Запись нового пользователя в БД
	err = u.Repo.Create(context.TODO(), newUser)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>Create", err)
		return nil, nil, ErrCreateUser
	}

	// Создаем токен по логину
	token, err := u.JWTServ.CreateToken(newUser.Login, newUser.Role)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>CreateToken", err)
		return nil, nil, ErrCreateUser
	}

	// Создаем Cookie
	cookie := u.Core.createCookie(token, u.Core.Args.GetNameCookie())

	// Формируем ответ
	userByte, err := json.Marshal(newUser)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>Marshal", err)
		return nil, nil, ErrCreateUser
	}

	return userByte, cookie, nil
}

func (u *Auth) Authentication(req *http.Request) ([]byte, *http.Cookie, error) {
	// Достаем login, password
	login, password, err := u.Core.takeLoginAndPassword(req)
	if err != nil {
		return nil, nil, ErrLoginPasswordIsBad2
	}

	// Идем в БД
	user, err := u.Repo.Read(context.TODO(), login)
	if err != nil {
		return nil, nil, ErrLogindNotFound
	}

	// Сверка паролей
	if !pkg.CompareHashAndPassword(user.Password, password) {
		return nil, nil, ErrPasswordNotValid
	}

	// Выдаем новый токен
	token, err := u.JWTServ.CreateToken(user.Login, user.Role)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Authentication>CreateToken", err)
		return nil, nil, err
	}

	return MessWelcome, u.Core.createCookie(token, u.Core.Args.GetNameCookie()), nil
}
