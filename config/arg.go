package config

import (
	"os"
	"strings"
)

type Flags struct {
	DryRun bool
}

var (
	args  []string
	flags Flags
)

func ParseArgs() {
	args := os.Args[1:]
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			// not a flag
			args = append(args, arg)
		}

		if arg == "--dry-run" {
			flags.DryRun = true
		}
	}
}

func GetFlags() Flags {
	return flags
}

func GetArgs() []string {
	return args
}
