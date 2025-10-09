package auth

import (
	"fmt"
	"time"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/golang-jwt/jwt/v4"
)

// Тип утверждения токена
type Claims struct {
	jwt.RegisteredClaims
	Login string
	Role  string
}

type JWTServ struct {
	Args config.Argser
	Logg logg.Logger
}

func NewJWTServ(argser config.Argser, logger logg.Logger) *JWTServ {
	return &JWTServ{
		Args: argser,
		Logg: logger,
	}
}

func (j *JWTServ) CreateToken(login, role string) (fullToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				// Settings of JWT
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Args.GetTimeLiveToken())), // Токен истекает через N сек
				NotBefore: jwt.NewNumericDate(time.Now()),                                // Токен активен с текущего момента
			},
			Login: login,
			Role:  role,
		})

	// Полный подписанный токен fullToken
	return token.SignedString(j.Args.GetSecretToken())
}

func (j *JWTServ) CheckOfValidToken(fullToken string) (login, role string, err error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(fullToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.Args.GetSecretToken(), nil
	})

	// Ошибка при проверке токена
	if err != nil || !token.Valid {
		return "", "", ErrTokenNotValid
	}

	// Проверка срока действия токена
	expirationTime := claims.ExpiresAt.Time
	if expirationTime.Before(time.Now().UTC()) {
		return "", "", ErrTokenIsBad
	}

	return claims.Login, claims.Role, nil
}
