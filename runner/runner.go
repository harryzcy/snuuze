package runner

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/runner/manager"
	"github.com/harryzcy/snuuze/runner/updater"
)

func RunForRepo(gitURL string) error {
	repoPath, err := prepareRepo(gitURL)
	if err != nil {
		return err
	}
	defer cleanupRepo(repoPath)

	infos, err := manager.Run(gitURL, repoPath)
	if err != nil {
		return err
	}

	updater.Update(gitURL, repoPath, infos)
	return nil
}

func prepareRepo(gitURL string) (gitPath string, err error) {
	cliConfig := config.GetCLIConfig()
	if !cliConfig.InPlace {
		gitPath, err = cloneRepo(gitURL)
		if err != nil {
			return "", err
		}

		err = updateGitCommitter(gitURL, gitPath)
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

	err := removeRepo(path)
	if err != nil {
		fmt.Println("Failed to remove repo:", err)
	}
}
