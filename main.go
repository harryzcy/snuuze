package main

import (
	"fmt"
	"log"

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

	if config.GetArgs().AsServer() {
		server.Init()
	} else {
		runCli()
	}
}

func runCli() {
	gitURL, repoPath, err := prepareRepo()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cleanupRepo(repoPath)

	manager.Run(gitURL, repoPath)
}

func prepareRepo() (gitURL, gitPath string, err error) {
	args := config.GetArgs()
	if len(args) != 0 {
		gitURL = args[0]
	} else {
		var err error
		gitURL, err = gitutil.GetOriginURL()
		if err != nil {
			return "", "", err
		}
	}

	flags := config.GetFlags()
	if !flags.InPlace {
		gitPath, err = gitutil.CloneRepo(gitURL)
		if err != nil {
			return "", "", err
		}

		err = gitutil.UpdateCommitter(gitURL, gitPath)
		if err != nil {
			return "", "", err
		}
	} else {
		gitPath = "."
	}

	return gitURL, gitPath, nil
}

func cleanupRepo(path string) {
	flags := config.GetFlags()
	if flags.InPlace {
		return
	}

	err := gitutil.RemoveRepo(path)
	if err != nil {
		fmt.Println("Failed to remove repo:", err)
	}
}
