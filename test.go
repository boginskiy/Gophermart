package main

import (
	"fmt"

	"github.com/boginskiy/Gophermart/pkg"
)

func main() {

	l := pkg.NewLuna()
	orderID := l.GenDigits(8)

	fmt.Println(orderID)

}
