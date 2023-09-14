package runner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harryzcy/snuuze/command"
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform"
)

var (
	ErrNoUserID = errors.New("app user ID is not set")
)

// GetGitOriginURL returns url of `origin` remote of the current git repo
func GetGitOriginURL() (string, error) {
	output, err := command.RunCommand(command.CommandInputs{
		Command: []string{"git", "remote", "get-url", "origin"},
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.Stdout.String()), nil
}

// cloneRepo clones a git repo to a temp directory
func cloneRepo(gitURL string) (string, error) {
	dirPath, err := os.MkdirTemp(config.TempDir(), "snuuze-*")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cloning repo to", dirPath)
	_, err = command.RunCommand(command.CommandInputs{
		Command: []string{"git", "clone", gitURL, dirPath},
	})
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func updateGitCommitter(gitURL, dirPath string) error {
	// TODO: support other git platforms
	if gitPlatform, _ := platform.DetermineGitPlatform(gitURL); gitPlatform != platform.GitPlatformGitHub {
		return nil
	}

	appName := config.GetHostingConfig().GitHub.AppName
	if appName == "" {
		appName = "snuuze"
	}
	if !strings.HasSuffix(appName, "[bot]") {
		appName = appName + "[bot]"
	}

	appUserID := config.GetHostingConfig().GitHub.AppUserID
	if appUserID == 0 {
		return ErrNoUserID
	}

	_, err := command.RunCommand(command.CommandInputs{
		Dir:     dirPath,
		Command: []string{"git", "config", "user.name", appName},
	})
	if err != nil {
		return err
	}

	_, err = command.RunCommand(command.CommandInputs{
		Dir:     dirPath,
		Command: []string{"git", "config", "user.email", fmt.Sprintf("%d+%s@users.noreply.github.com", appUserID, appName)},
	})
	if err != nil {
		return err
	}

	_, err = command.RunCommand(command.CommandInputs{
		Dir:     dirPath,
		Command: []string{"git", "config", "commit.gpgsign", "false"},
	})
	if err != nil {
		return err
	}

	return nil
}

func removeRepo(path string) error {
	return os.RemoveAll(path)
}
