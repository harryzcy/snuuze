package config

import (
	"fmt"
	"os"
	"strings"
)

type CLIConfig struct {
	Mode    string
	Args    []string
	DryRun  bool
	InPlace bool
}

func (c CLIConfig) AsServer() bool {
	return c.Mode == "server"
}

var (
	cliConfig CLIConfig
)

func ParseArgs() {
	args := os.Args[1:]
	if len(args) == 0 {
		cliConfig.Mode = "server"
	} else if args[0] == "server" {
		cliConfig.Mode = "server"
		args = args[1:]
	} else if args[0] == "cli" {
		cliConfig.Mode = "cli"
		args = args[1:]
	} else {
		fmt.Fprintln(os.Stderr, "error: invalid mode")
		os.Exit(1)
	}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			// not a flag
			cliConfig.Args = append(cliConfig.Args, arg)
		}

		if arg == "--dry-run" {
			cliConfig.DryRun = true
		} else if arg == "--in-place" {
			cliConfig.InPlace = true
		} else {
			fmt.Fprintln(os.Stderr, "error: invalid flag: "+arg)
			os.Exit(1)
		}
	}

	if cliConfig.InPlace && !cliConfig.DryRun {
		fmt.Fprintln(os.Stderr, "error: --in-place must be used together with --dry-run")
		os.Exit(1)
	}
}

func GetCLIConfig() CLIConfig {
	return cliConfig
}
