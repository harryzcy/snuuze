package main

import (
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
		log.Fatal(err)
		return
	}
	err = config.LoadHostingConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	if config.GetCLIConfig().AsServer() {
		runServer()
	} else {
		runCli()
	}
}

func runServer() {
	server.Run()
}

func runCli() {
	var gitURL string
	cliConfig := config.GetCLIConfig()
	if len(cliConfig.Args) != 0 {
		gitURL = cliConfig.Args[0]
	} else {
		var err error
		gitURL, err = git.GetOriginURL()
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	err := runner.RunForRepo(gitURL)
	if err != nil {
		log.Fatal(err)
	}
}
