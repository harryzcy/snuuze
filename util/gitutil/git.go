package gitutil

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/util/cmdutil"
)

var (
	ErrNoUserID = errors.New("app user ID is not set")
)

func GetOriginURL() (string, error) {
	output, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "remote", "get-url", "origin"},
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.Stdout.String()), nil
}

// CloneRepo clones a git repo to a temp directory
func CloneRepo(gitURL string) (string, error) {
	dirPath, err := os.MkdirTemp(config.TempDir(), "snuuze-*")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cloning repo to", dirPath)
	_, err = cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"git", "clone", gitURL, dirPath},
	})
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func UpdateCommitter(gitURL, dirPath string) error {
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

	_, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Dir:     dirPath,
		Command: []string{"git", "config", "user.name", appName},
	})
	if err != nil {
		return err
	}

	_, err = cmdutil.RunCommand(cmdutil.CommandInputs{
		Dir:     dirPath,
		Command: []string{"git", "config", "user.email", fmt.Sprintf("%d+%s@users.noreply.github.com", appUserID, appName)},
	})
	if err != nil {
		return err
	}

	_, err = cmdutil.RunCommand(cmdutil.CommandInputs{
		Dir:     dirPath,
		Command: []string{"git", "config", "commit.gpgsign", "false"},
	})
	if err != nil {
		return err
	}

	return nil
}

func RemoveRepo(path string) error {
	return os.RemoveAll(path)
}