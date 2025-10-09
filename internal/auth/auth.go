package auth

import (
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/models"
	"github.com/boginskiy/Gophermart/pkg"
)

type Auth struct {
	Repo    repository.Repository
	JWTServ JWTokener
	Core    *AhCore
}

func NewAuth(core *AhCore, jwter JWTokener, repo repository.Repository) *Auth {
	return &Auth{
		JWTServ: jwter,
		Repo:    repo,
		Core:    core,
	}
}

func (u *Auth) CheckToken(token string) (login, role string, err error) {
	return u.JWTServ.CheckOfValidToken(token)
}

func (u *Auth) GetCredentials(req *http.Request) (login, password string, err error) {
	// TODO!
	// Остановка тут >>
	// Как можно быстрее определеитьчто массив байт содержит необходимые данные
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

	// Создаем токен по логину
	token, err := u.JWTServ.CreateToken(newUser.Login, newUser.Role)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Registration>CreateToken", err)
		return nil, nil, ErrCreateUser
	}

	// Создаем Cookie
	cookie := u.Core.CreateCookie(token, u.Core.Args.GetNameCookie())

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
		return nil, nil, err
	}

	// Идем в БД
	recore, err := u.Repo.SelectRecord(login)
	if err != nil {
		// TODO! Пока ошибка только, что login неверный,
		// потом, возможно стоит расширить спектр. Например 500 ошибка
		return nil, nil, ErrLogindNotFound
	}

	user, _ := recore.(*models.User)

	// Сверка паролей входящего и пароля из БД
	inputPassword, err := pkg.GenerateHash(password)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Authentication>GenerateHash", err)
		return nil, nil, err
	}

	// Пароль не прошел сверку с предыдущим паролем
	if !pkg.CompareHashAndPassword(user.Password, inputPassword) {
		return nil, nil, ErrPasswordNotValid
	}

	// Выдаем новый токен
	token, err := u.JWTServ.CreateToken(user.Login, user.Role)
	if err != nil {
		u.Core.Logg.RaiseError("Auth>Authentication>CreateToken", err)
		return nil, nil, err
	}

	return MessWelcome, u.Core.CreateCookie(token, u.Core.Args.GetNameCookie()), nil
}
