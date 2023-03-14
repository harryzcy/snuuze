package main

import (
	"log"
	"os"

	"github.com/harryzcy/snuuze/checker"
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/gitutil"
	"github.com/harryzcy/snuuze/matcher"
	"github.com/harryzcy/snuuze/updater"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	gitURL, repoPath, err := prepareRepo()
	if err != nil {
		log.Fatal(err)
		return
	}

	matches, err := matcher.Scan(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	infos, err := checker.ListUpgrades(matches)
	if err != nil {
		log.Fatal(err)
	}

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
