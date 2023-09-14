package runner

import (
	"fmt"
	"log"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/manager"
	"github.com/harryzcy/snuuze/util/gitutil"
)

func RunForRepo(gitURL string) {
	cliConfig := config.GetCLIConfig()
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
