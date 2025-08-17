package hostsys

import (
	"bytes"
	"os/exec"
)

type Executor struct {
}

type ExecOpt func(*exec.Cmd) error

func ChdirOpt(dir string) ExecOpt {
	return func(cmd *exec.Cmd) error {
		cmd.Dir = dir

		return nil
	}
}

func (e *Executor) Exec(cmdName string, args []string, opts ...ExecOpt) (stdout string, stderr string, err error) {
	cmd := exec.Command(cmdName, args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	for _, opt := range opts {
		opt(cmd)
	}

	err = cmd.Run()

	return outb.String(), errb.String(), err
}

func NewExecutor() *Executor {
	return &Executor{}
}
