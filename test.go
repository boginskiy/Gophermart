package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - собственное утверждение
type Claims struct {
	jwt.RegisteredClaims
	Login string
	Role  string
}

func CreateToken(login, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				// Settings of JWT
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Second)), // Токен истекает через N
				NotBefore: jwt.NewNumericDate(time.Now()),                       // Токен активен с текущего момента
			},
			Login: login,
			Role:  role,
		})

	// Полный подписанный токен
	return token.SignedString([]byte("SecretKey"))
}

func main() {

}
