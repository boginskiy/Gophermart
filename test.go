package main

import (
	"fmt"

	"github.com/boginskiy/Gophermart/pkg"
)

func main() {

	password := "ЛОХ777"

	// hashByte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// fmt.Println(hashByte)
	// err := bcrypt.CompareHashAndPassword(hashByte, []byte(password))
	// fmt.Println(err == nil)

	pass, _ := pkg.GenerateHash(password)

	res := pkg.CompareHashAndPassword(pass, password)

	fmt.Println(res)

}
