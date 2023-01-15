package main

import (
	"fmt"

	"github.com/harryzcy/sailor/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}
