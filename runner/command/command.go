package command

import (
	"bytes"
	"errors"
	"os/exec"
	"time"

	"github.com/harryzcy/snuuze/config"
)

var (
	ErrorTimeout = errors.New("timeout")
)

type Inputs struct {
	Command []string
	Dir     string
	Env     map[string]string
	Timeout time.Duration
}

func (input *Inputs) GetTimeout() time.Duration {
	if input.Timeout <= 0 {
		return config.GetHostingConfig().Data.GetTimeout()
	}
	return input.Timeout
}

type Output struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func RunCommand(inputs Inputs) (*Output, error) {
	cmd := exec.Command(inputs.Command[0], inputs.Command[1:]...) // #nosec G204
	cmd.Dir = inputs.Dir
	cmd.Env = []string{}
	for key, value := range inputs.Env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// set timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	select {
	case <-time.After(inputs.GetTimeout()):
		if err := cmd.Process.Kill(); err != nil {
			return nil, err
		}
		return nil, ErrorTimeout
	case err := <-done:
		return &Output{
			Stdout: out,
			Stderr: stderr,
		}, err
	}
}
