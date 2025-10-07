package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
)

type AhCore struct {
	Args config.Argser
	Logg logg.Logger
}

func NewAhCore(argser config.Argser, logger logg.Logger) *AhCore {
	return &AhCore{Args: argser, Logg: logger}
}

func (c *AhCore) takeLoginAndPassword(req *http.Request) (login, password string, err error) {
	// Читаем body req
	dataByte, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		c.Logg.RaiseError("Auth>takeLoginAndPassword>ReadAll", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Парсинг данных
	overallStruct := make(map[string]any, 2)
	err = json.Unmarshal(dataByte, &overallStruct)
	if err != nil {
		c.Logg.RaiseError("Auth>takeLoginAndPassword>Unmarshal", err)
		return "", "", ErrLoginPasswordIsBad
	}

	// Достаем login, password
	if login, ok := overallStruct["login"].(string); ok {
		if password, ok2 := overallStruct["password"].(string); ok2 {
			return login, password, nil
		}
	}

	c.Logg.RaiseInfo("Auth>takeLoginAndPassword>data not found")
	return "", "", ErrLoginPasswordIsBad
}

func (c *AhCore) CreateCookie(token, name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // Доступ только серверу, увеличивает безопасность
		SameSite: http.SameSiteStrictMode, // Запрещает отправлять куки с другого домена
		MaxAge:   300,                     // Жива 300 секунд
		Secure:   false,                   // Поставьте true, если работаете через HTTPS
	}
}
