package main

import (
	"fmt"

	"github.com/harryzcy/sailor/config"
	"github.com/harryzcy/sailor/manager"
)

func main() {
	err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := manager.Scan()
	fmt.Println(files)
}
