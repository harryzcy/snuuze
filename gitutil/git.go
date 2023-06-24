package gitutil

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harryzcy/snuuze/cmdutil"
)

var (
	TEMP_DIR = os.Getenv("TEMP_DIR")
)

func init() {
	if TEMP_DIR == "" {
		TEMP_DIR = os.TempDir()
	}
}

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
	dirPath, err := os.MkdirTemp(TEMP_DIR, "snuuze-*")
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

func RemoveRepo(path string) error {
	return os.RemoveAll(path)
}
