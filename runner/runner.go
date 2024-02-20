package runner

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/runner/git"
	"github.com/harryzcy/snuuze/runner/manager"
	"github.com/harryzcy/snuuze/runner/updater"
	"github.com/harryzcy/snuuze/types"
)

func RunForRepo(gitURL string) error {
	cliConfig := config.GetCLIConfig()

	repoPath, err := prepareRepo(gitURL)
	if err != nil {
		return err
	}
	defer cleanupRepo(repoPath)

	infos, err := manager.Run(repoPath)
	if err != nil {
		return err
	}

	if cliConfig.InPlace {
		manager.PrintUpgradeInfos(infos)
	}

	err = updater.Update(gitURL, repoPath, infos, !cliConfig.InPlace)
	return err
}

func GetDependencyForRepo(gitURL string) (map[types.PackageManager][]*types.Dependency, error) {
	repoPath, err := prepareRepo(gitURL)
	if err != nil {
		return nil, err
	}
	defer cleanupRepo(repoPath)

	dependencies, err := manager.FindAll(repoPath)
	if err != nil {
		return nil, err
	}
	return dependencies, nil
}

func prepareRepo(gitURL string) (gitPath string, err error) {
	cliConfig := config.GetCLIConfig()
	if !cliConfig.InPlace {
		gitPath, err = git.CloneRepo(gitURL)
		if err != nil {
			return "", err
		}

		err = git.UpdateCommitter(gitURL, gitPath)
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

	err := git.RemoveRepo(path)
	if err != nil {
		fmt.Println("Failed to remove repo:", err)
	}
}
