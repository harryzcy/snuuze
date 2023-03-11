package cmdutil

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		command []string
		env     map[string]string
		stdout  string
		stderr  string
	}{
		{
			command: []string{"echo", "hello"},
			stdout:  "hello\n",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			output, err := RunCommand(CommandInputs{
				Command: test.command,
				Env:     test.env,
			})
			assert.Nil(t, err)
			assert.Equal(t, test.stdout, output.Stdout.String())
			assert.Equal(t, test.stderr, output.Stderr.String())
		})
	}
}
