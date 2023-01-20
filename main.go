package main

import (
	"fmt"
	"log"

	"github.com/harryzcy/latte/checker"
	"github.com/harryzcy/latte/config"
	"github.com/harryzcy/latte/matcher"
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

	infos, err := checker.ListUpgrades(matches)
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range infos {
		fmt.Println(info)
	}
}
