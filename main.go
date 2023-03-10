package main

import (
	"fmt"
	"log"

	"github.com/harryzcy/snuuze/checker"
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/matcher"
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
