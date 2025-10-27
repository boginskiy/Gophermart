package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	repo "github.com/boginskiy/Gophermart/internal/repository"
	"github.com/boginskiy/Gophermart/models"
	"github.com/boginskiy/Gophermart/pkg"
)

type Auth struct {
	Config  conf.Config
	Logger  logg.Logger
	Repo    repo.RepoUsersTber
	JWTServ JWTokener
}

func NewAuth(config conf.Config, logger logg.Logger, repoUsers repo.RepoUsersTber, jwter JWTokener) *Auth {
	return &Auth{
		Config:  config,
		Logger:  logger,
		Repo:    repoUsers,
		JWTServ: jwter,
	}
}

func (u *Auth) createCookie(token, name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                         // Доступ только серверу, увеличивает безопасность
		SameSite: http.SameSiteStrictMode,      // Запрещает отправлять куки с другого домена
		MaxAge:   u.Config.GetTimeLiveCookie(), // Жива N секунд
		Secure:   false,                        // Поставьте true, если работаете через HTTPS
	}
}

func (u *Auth) takeLoginAndPassword(req *http.Request) (login, password string, err error) {
	// Читаем body req
	dataByte, err := io.ReadAll(req.Body)

	defer req.Body.Close()
	if err != nil {
		u.Logger.RaiseError("Auth>takeLoginAndPassword>ReadAll", err)
		return "", "", err
	}

	// Парсинг данных
	var overallStruct = make(map[string]any)
	err = json.Unmarshal(dataByte, &overallStruct)
	if err != nil {
		u.Logger.RaiseError("Auth>takeLoginAndPassword>Unmarshal", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Достаем login, password
	if login, ok := overallStruct["login"].(string); ok {
		if password, ok2 := overallStruct["password"].(string); ok2 {
			return login, password, nil
		}
	}

	u.Logger.RaiseInfo("Auth>takeLoginAndPassword>data not found")
	return "", "", ErrLoginPasswordIsBad
}

func (u *Auth) CheckToken(token string) (login, role string, id int64, err error) {
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
	login, password, err := u.takeLoginAndPassword(req)
	if err != nil {
		return nil, nil, ErrLoginPasswordIsBad
	}

	// Проверка уникальности login
	loginIsUnic, _ := u.Repo.CheckUnic(context.TODO(), login)
	if !loginIsUnic {
		// Пользователь не уникален, логин уже занят
		u.Logger.RaiseInfo(ErrLoginNotUnic.Error())
		return nil, nil, ErrLoginNotUnic
	}

	// Создаем нового пользователя
	newUser, err := models.NewUser(login, password)
	if err != nil {
		u.Logger.RaiseError("Auth>Registration>NewUser", err)
		return nil, nil, ErrCreateUser
	}

	// Запись нового пользователя в БД
	newUserID, err := u.Repo.Create(context.TODO(), newUser)
	if err != nil {
		u.Logger.RaiseError("Auth>Registration>Create", err)
		return nil, nil, ErrCreateUser
	}

	// Создаем токен по логину
	token, err := u.JWTServ.CreateToken(newUser.Login, newUser.Role, newUserID)
	if err != nil {
		u.Logger.RaiseError("Auth>Registration>CreateToken", err)
		return nil, nil, ErrCreateUser
	}

	// Создаем Cookie
	cookie := u.createCookie(token, u.Config.GetNameCookie())

	// Формируем ответ
	userByte, err := json.Marshal(newUser)
	if err != nil {
		u.Logger.RaiseError("Auth>Registration>Marshal", err)
		return nil, nil, ErrCreateUser
	}

	return userByte, cookie, nil
}

func (u *Auth) Authentication(req *http.Request) ([]byte, *http.Cookie, error) {
	// Достаем login, password
	login, password, err := u.takeLoginAndPassword(req)
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
	token, err := u.JWTServ.CreateToken(user.Login, user.Role, user.ID)
	if err != nil {
		u.Logger.RaiseError("Auth>Authentication>CreateToken", err)
		return nil, nil, err
	}

	return MessWelcome, u.createCookie(token, u.Config.GetNameCookie()), nil
}
