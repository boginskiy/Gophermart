package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type AhCore struct {
	Args    config.Argser
	Logg    logg.Logger
	JWTServ JWTokener
}

func NewAhCore(argser config.Argser, logger logg.Logger) *AhCore {
	return &AhCore{Args: argser, Logg: logger}
}

func (c *AhCore) CreateCookie(token, name string) *http.Cookie {
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

func (c *AhCore) takeLoginAndPassword(req *http.Request) (login, password string, err error) {
	// Читаем body req
	dataByte, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		c.Logg.RaiseError("AhCore>takeLoginAndPassword>ReadAll", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Парсинг данных
	overallStruct := make(map[string]any, 2)
	err = json.Unmarshal(dataByte, &overallStruct)
	if err != nil {
		c.Logg.RaiseError("AhCore>takeLoginAndPassword>Unmarshal", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Достаем login, password
	if login, ok := overallStruct["login"].(string); ok {
		if password, ok2 := overallStruct["password"].(string); ok2 {
			return login, password, nil
		}
	}

	c.Logg.RaiseInfo("AhCore>takeLoginAndPassword>data not found")
	return "", "", ErrLoginPasswordIsBad
}

// var CtxLogin = CtxKey{}
// var CtxRole = CtxKey{}
func (c *AhCore) takeParamFromCtx(req *http.Request, p CtxKey) string {
	param, ok := req.Context().Value(p).(string)
	if !ok {
		c.Logg.RaiseError("AhCore>takeParamFromCtx", nil)
		return ""
	}
	return param
}
