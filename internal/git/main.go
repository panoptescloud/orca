package git

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
)

const ConfirmYes = "Yes"
const ConfirmNo = "No"

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
	RecordIfError(msg string, err error) error
	NewLine()
	PresentChoices(opts []string, title string) (string, error)
}

type executor interface {
	Exec(cmdName string, args []string, opts ...hostsys.ExecOpt) error
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
	err := g.exec.Exec("git", []string{"rev-parse", "--is-inside-work-tree"})

	return err == nil
}

func (g *Git) GetRepositoryRootFromPath(path string) (string, error) {
	stdoutOpt, stdout := hostsys.WithStdout()
	stderrOpt, stderr := hostsys.WithStderr()

	dir := filepath.Dir(path)

	err := g.exec.Exec("git", []string{
		"rev-parse",
		"--show-toplevel",
	}, hostsys.ChdirOpt(dir), stdoutOpt, stderrOpt)

	if err != nil {
		errOutput := stderr.String()
		slog.Error("failed to execute 'git rev-parse --show-top-level'", "in", dir, "err", err, "stderr", errOutput)

		return "", common.ErrCommandExecutionFailed{
			Msg: errOutput,
		}
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (g *Git) DirectoryHasOrigin(dir string, origin string) (bool, error) {
	stdoutOpt, stdout := hostsys.WithStdout()

	err := g.exec.Exec("git", []string{
		"remote",
		"get-url",
		"origin",
	}, hostsys.ChdirOpt(dir), stdoutOpt)

	if err != nil {
		return false, err
	}

	return strings.TrimSpace(stdout.String()) == origin, nil
}

func NewGit(exec executor, tui tui) *Git {
	return &Git{
		exec: exec,
		tui:  tui,
	}
}
