package workspaces

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/spf13/afero"
)

const DefaultWorkspaceFileName = "orca.workspace.yaml"

type tui interface {
	Info(msg ...string)
	Error(msg ...string)
	Success(msg ...string)
	RecordIfError(msg string, err error) error
}

type configManager interface {
	AddWorkspace(path string, ws *common.Workspace) error
	GetWorkspaceLocations() common.WorkspaceLocations
}

type Manager struct {
	fs            afero.Fs
	tui           tui
	configManager configManager
}

func NewManager(fs afero.Fs, tui tui, configManager configManager) *Manager {
	return &Manager{
		fs:            fs,
		tui:           tui,
		configManager: configManager,
	}
}
