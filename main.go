package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harryzcy/snuuze/checker"
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/gitutil"
	"github.com/harryzcy/snuuze/matcher"
	"github.com/harryzcy/snuuze/updater"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	gitURL, repoPath, err := prepareRepo()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cleanupRepo(repoPath)

	matches, err := matcher.Scan(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	infos, err := checker.ListUpgrades(matches)
	if err != nil {
		log.Fatal(err)
	}

	if len(infos) == 0 {
		fmt.Println("No updates found")
		return
	}
	fmt.Println("Found", len(infos), "updates")

	updater.Update(gitURL, repoPath, infos)
}

func prepareRepo() (gitURL, path string, err error) {
	if len(os.Args) == 2 {
		gitURL = os.Args[1]
	} else {
		var err error
		gitURL, err = gitutil.GetOriginURL()
		if err != nil {
			return "", "", err
		}
	}

	path, err = gitutil.CloneRepo(gitURL)
	if err != nil {
		return "", "", err
	}
	return gitURL, path, nil
}

func cleanupRepo(path string) {
	gitutil.RemoveRepo(path)
}
