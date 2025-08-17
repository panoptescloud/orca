package git

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
	RecordIfError(msg string, err error) error
	NewLine()
	PresentChoices(opts []string, title string) (string, error)
}

type executor interface {
	Exec(cmdName string, args []string, opts ...hostsys.ExecOpt) (string, string, error)
}

type Git struct {
	exec executor
	tui  tui
}

func (g *Git) mustBeInAGitRepository() error {
	if g.isInGitRepository() {
		return nil
	}

	return g.tui.RecordIfError(
		"You don't appear to be in a git repository!",
		common.ErrInvalidExecutionContext{
			Msg: "this command must be executed from within a git repository",
		},
	)

}

func (g *Git) isInGitRepository() bool {
	_, _, err := g.exec.Exec("git", []string{"rev-parse", "--is-inside-work-tree"})

	return err == nil
}

func NewGit(exec executor, tui tui) *Git {
	return &Git{
		exec: exec,
		tui:  tui,
	}
}
