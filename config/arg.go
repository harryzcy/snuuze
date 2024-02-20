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

const (
	ModeServer = "server"
	ModeCLI    = "cli"
)

func (c CLIConfig) AsServer() bool {
	return c.Mode == ModeServer
}

var (
	cliConfig CLIConfig
)

func ParseArgs() {
	args := os.Args[1:]
	switch {
	case len(args) == 0:
		cliConfig.Mode = ModeServer
	case args[0] == "server":
		cliConfig.Mode = ModeServer
		args = args[1:]
	case args[0] == "cli":
		cliConfig.Mode = ModeCLI
		args = args[1:]
	default:
		fmt.Fprintln(os.Stderr, "error: invalid mode")
		os.Exit(1)
	}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			// not a flag
			cliConfig.Args = append(cliConfig.Args, arg)
		}

		switch arg {
		case "--in-place":
			cliConfig.InPlace = true
		case "--dry-run":
			cliConfig.DryRun = true
		default:
			fmt.Fprintln(os.Stderr, "error: invalid flag: "+arg)
			os.Exit(1)
		}
	}
}

func GetCLIConfig() CLIConfig {
	return cliConfig
}
