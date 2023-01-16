package main

import (
	"fmt"
	"log"

	"github.com/harryzcy/sailor/checker"
	"github.com/harryzcy/sailor/config"
	"github.com/harryzcy/sailor/matcher"
)

func main() {
	err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	matches, err := matcher.Scan()
	if err != nil {
		log.Fatal(err)
	}

	checker.ListUpgrades(matches)
}
