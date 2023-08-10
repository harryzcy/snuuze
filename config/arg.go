package config

import (
	"fmt"
	"os"
	"strings"
)

type Flags struct {
	DryRun  bool
	InPlace bool
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

		if arg == "--in-place" {
			flags.InPlace = true
		}
	}

	if flags.DryRun && flags.InPlace {
		fmt.Fprintln(os.Stderr, "cannot have both --dry-run and --in-place")
		os.Exit(1)
	}
}

func GetFlags() Flags {
	return flags
}

func GetArgs() []string {
	return args
}
