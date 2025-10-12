package pkg

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Luna struct {
}

func NewLuna() *Luna {
	return &Luna{}
}

func (l *Luna) genRandomDigitsStr(long int) string {
	var result string
	for i := 0; i < long; i++ {
		result += fmt.Sprintf("%d", rand.Intn(10))
	}
	return result
}

func (l *Luna) preparCheckSum(digits string) string {
	sum := 0
	isDouble := true

	// Обход справа налево
	for i := len(digits) - 1; i >= 0; i-- {
		digit := int(digits[i] - '0')

		if isDouble {
			digit *= 2

			// Если удвоенное число больше 9, уменьшаем на 9
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isDouble = !isDouble
	}
	return strconv.Itoa(10 - (sum%10)%10)
}

func (l *Luna) CheckDigits(digits string) bool {
	sum := 0
	isDouble := false

	// Обход справа налево
	for i := len(digits) - 1; i >= 0; i-- {
		digit := int(digits[i] - '0')

		if isDouble {
			digit *= 2

			// Если удвоенное число больше 9, уменьшаем на 9
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isDouble = !isDouble
	}
	return sum%10 == 0
}

func (l *Luna) GenDigits(long int) string {
	setDigits := l.genRandomDigitsStr(long - 1)
	endDigit := l.preparCheckSum(setDigits)
	return setDigits + endDigit
}
