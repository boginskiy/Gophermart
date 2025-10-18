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

// 96357934
// 14884589
// 62028444
// 35953629
// 68379148
