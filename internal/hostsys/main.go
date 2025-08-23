package hostsys

import "github.com/spf13/afero"

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
	RecordIfError(msg string, err error) error
}

type HostSystem struct {
	tui      tui
	fs       afero.Fs
	toolsDir string
}

func NewHostSystem(tui tui, fs afero.Fs, toolsDir string) *HostSystem {
	return &HostSystem{
		tui:      tui,
		fs:       fs,
		toolsDir: toolsDir,
	}
}
