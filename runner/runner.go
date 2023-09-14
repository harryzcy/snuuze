package runner

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/manager"
	"github.com/harryzcy/snuuze/util/gitutil"
)

func RunForRepo(gitURL string) error {
	repoPath, err := prepareRepo(gitURL)
	if err != nil {
		return err
	}
	defer cleanupRepo(repoPath)

	return manager.Run(gitURL, repoPath)
}

func prepareRepo(gitURL string) (gitPath string, err error) {
	cliConfig := config.GetCLIConfig()
	if !cliConfig.InPlace {
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
