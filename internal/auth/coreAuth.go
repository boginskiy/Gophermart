package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type CoreAh struct {
	Args    config.Argser
	Logg    logg.Logger
	JWTServ JWTokener
}

func NewCoreAh(argser config.Argser, logger logg.Logger) *CoreAh {
	return &CoreAh{Args: argser, Logg: logger}
}

func (c *CoreAh) createCookie(token, name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                       // Доступ только серверу, увеличивает безопасность
		SameSite: http.SameSiteStrictMode,    // Запрещает отправлять куки с другого домена
		MaxAge:   c.Args.GetTimeLiveCookie(), // Жива N секунд
		Secure:   false,                      // Поставьте true, если работаете через HTTPS
	}
}

func (c *CoreAh) takeLoginAndPassword(req *http.Request) (login, password string, err error) {
	// Читаем body req
	dataByte, err := io.ReadAll(req.Body)

	defer req.Body.Close()
	if err != nil {
		c.Logg.RaiseError("CoreAh>takeLoginAndPassword>ReadAll", err)
		return "", "", err
	}

	// Парсинг данных
	var overallStruct = make(map[string]any)
	err = json.Unmarshal(dataByte, &overallStruct)
	if err != nil {
		c.Logg.RaiseError("CoreAh>takeLoginAndPassword>Unmarshal", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Достаем login, password
	if login, ok := overallStruct["login"].(string); ok {
		if password, ok2 := overallStruct["password"].(string); ok2 {
			return login, password, nil
		}
	}

	c.Logg.RaiseInfo("CoreAh>takeLoginAndPassword>data not found")
	return "", "", ErrLoginPasswordIsBad
}

func (c *CoreAh) takeParamFromCtx(req *http.Request, p CtxKey) string {
	param, ok := req.Context().Value(p).(string)
	if !ok {
		c.Logg.RaiseError("CoreAh>takeParamFromCtx", nil)
		return ""
	}
	return param
}
