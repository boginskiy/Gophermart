package main

import (
	"fmt"
	"time"
)

func main() {
	//
	// t := time.Now().Format(time.RFC3339)

	// fmt.Println(t)

	utcNow := time.Now().UTC()
	fmt.Println(utcNow, utcNow.Format(time.RFC3339))
}
