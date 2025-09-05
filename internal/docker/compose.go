package docker

import (
	"github.com/panoptescloud/orca/internal/hostsys"
)

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
	NewLine()
	RecordIfError(msg string, err error) error
}

type cli interface {
	Exec(cmdName string, args []string, opts ...hostsys.ExecOpt) (err error)
}

type Compose struct {
	cli cli
	tui tui
}

func NewCompose(cli cli, tui tui) *Compose {
	return &Compose{
		cli: cli,
		tui: tui,
	}
}
