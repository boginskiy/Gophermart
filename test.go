package main

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now()
	formattedTime := now.Format(time.RFC3339)

	fmt.Println(formattedTime)

	tt := time.Now()

	fmt.Println(tt.Format(time.RFC3339))

}
