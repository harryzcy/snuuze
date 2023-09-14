package main

import (
	"log"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/runner"
	"github.com/harryzcy/snuuze/server"
	"github.com/harryzcy/snuuze/util/gitutil"
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
		server.Run()
	} else {
		runCli()
	}
}

func runCli() {
	var gitURL string
	cliConfig := config.GetCLIConfig()
	if len(cliConfig.Args) != 0 {
		gitURL = cliConfig.Args[0]
	} else {
		var err error
		gitURL, err = gitutil.GetOriginURL()
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	runner.RunForRepo(gitURL)
}
