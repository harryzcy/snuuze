package main

import (
	"fmt"
	"log"

	"github.com/harryzcy/sailor/config"
	"github.com/harryzcy/sailor/matcher"
)

func main() {
	err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := matcher.Scan()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files)
}
