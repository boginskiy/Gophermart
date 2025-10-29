package auth

import (
	"fmt"
	"time"

	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/golang-jwt/jwt/v4"
)

// Тип утверждения токена
type Claims struct {
	jwt.RegisteredClaims
	Login string
	Role  string
	ID    int64
}

type JWTServ struct {
	Config conf.Config
	Logger logg.Logger
}

func NewJWTServ(config conf.Config, logger logg.Logger) *JWTServ {
	return &JWTServ{
		Config: config,
		Logger: logger,
	}
}

func (j *JWTServ) CreateToken(login, role string, id int64) (fullToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				// Settings of JWT
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Config.GetTimeLiveToken())), // Токен истекает через N сек
				NotBefore: jwt.NewNumericDate(time.Now()),                                  // Токен активен с текущего момента
			},
			Login: login,
			Role:  role,
			ID:    id,
		})

	// Полный подписанный токен fullToken
	return token.SignedString(j.Config.GetSecretKeyToken())
}

func (j *JWTServ) CheckOfValidToken(fullToken string) (login, role string, id int64, err error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(fullToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.Config.GetSecretKeyToken(), nil
	})

	// Ошибка при проверке токена
	if err != nil || !token.Valid {
		return "", "", 0, ErrTokenNotValid
	}

	// Проверка срока действия токена
	expirationTime := claims.ExpiresAt.Time
	if expirationTime.Before(time.Now().UTC()) {
		return "", "", 0, ErrTokenIsBad
	}

	return claims.Login, claims.Role, claims.ID, nil
}
