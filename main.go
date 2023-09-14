package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/manager"
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

	repoPath, err := prepareRepo(gitURL, cliConfig.InPlace)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cleanupRepo(repoPath)

	manager.Run(gitURL, repoPath)
}

func prepareRepo(gitURL string, inPlace bool) (gitPath string, err error) {
	if !inPlace {
		gitPath, err = gitutil.CloneRepo(gitURL)
		if err != nil {
			return "", err
		}

		err = gitutil.UpdateCommitter(gitURL, gitPath)
		if err != nil {
			return "", err
		}
	} else {
		gitPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	return gitPath, nil
}

func cleanupRepo(path string) {
	cliConfig := config.GetCLIConfig()
	if cliConfig.InPlace {
		return
	}

	err := gitutil.RemoveRepo(path)
	if err != nil {
		fmt.Println("Failed to remove repo:", err)
	}
}
