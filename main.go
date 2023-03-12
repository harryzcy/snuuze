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

	root, err := prepareRepo()
	if err != nil {
		log.Fatal(err)
		return
	}

	matches, err := matcher.Scan(root)
	if err != nil {
		log.Fatal(err)
	}

	infos, err := checker.ListUpgrades(matches)
	if err != nil {
		log.Fatal(err)
	}

	updater.Update(root, infos)
}

func prepareRepo() (string, error) {
	var gitUrl string
	if len(os.Args) == 2 {
		gitUrl = os.Args[1]
	} else {
		var err error
		gitUrl, err = gitutil.GetOriginURL()
		if err != nil {
			return "", err
		}
	}

	path, err := gitutil.CloneRepo(gitUrl)
	if err != nil {
		return "", err
	}
	return path, nil
}
