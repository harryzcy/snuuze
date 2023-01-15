package main

import (
	"fmt"
	"log"

	"github.com/harryzcy/sailor/config"
	"github.com/harryzcy/sailor/matcher"
	"github.com/harryzcy/sailor/parser"
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

	for _, match := range matches {
		dependencies, _ := parser.Parse(match)
		for _, dependency := range dependencies {
			fmt.Println(dependency)
		}
	}
}
