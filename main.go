package main

import (
	"fmt"
	"log"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/runner"
	"github.com/harryzcy/snuuze/runner/git"
	"github.com/harryzcy/snuuze/server"
)

func main() {
	config.ParseArgs()
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
		return
	}
	err = config.LoadHostingConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load hosting config: %w", err))
		return
	}

	if config.GetCLIConfig().AsServer() {
		runServer()
	} else {
		runCli()
	}
}

// runServer runs as a server
func runServer() {
	server.Run()
}

// runCli runs as a CLI tool
func runCli() {
	var gitURL string
	cliConfig := config.GetCLIConfig()
	if len(cliConfig.Args) != 0 {
		gitURL = cliConfig.Args[0]
	} else {
		var err error
		gitURL, err = git.GetOriginURL()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to get git origin url: %w", err))
			return
		}
	}

	err := runner.RunForRepo(gitURL)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to run for repo: %w", err))
	}
}
