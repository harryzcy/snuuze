package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/runner/command"
)

var (
	ErrNoUserID = errors.New("app user ID is not set")
)

// GetGitOriginURL returns url of `origin` remote of the current git repo
func GetOriginURL() (string, error) {
	output, err := command.RunCommand(command.CommandInputs{
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
		return "", err
	}

	fmt.Println("Cloning repo to", dirPath)

	gitPlatform, url := platform.DetermineGitPlatform(gitURL)
	client, err := platform.NewClient(platform.NewClientOptions{
		URL:      url,
		Platform: gitPlatform,
	})
	if err != nil {
		return "", err
	}
	authGitURL, err := getGitURLWithToken(client, gitURL)
	if err != nil {
		return "", err
	}

	_, err = command.RunCommand(command.CommandInputs{
		Command: []string{"git", "clone", authGitURL, dirPath},
	})
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

// getGitURLWithToken returns a git url with token set
func getGitURLWithToken(client platform.Client, gitURL string) (string, error) {
	ctx := context.Background()
	token, err := client.Token(ctx)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(gitURL, "https://", fmt.Sprintf("https://x-oauth-basic:%s@", token)), nil
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
		appName += "[bot]"
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

func RemoveRepo(path string) error {
	return os.RemoveAll(path)
}
