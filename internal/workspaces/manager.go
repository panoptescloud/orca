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

type config interface {
	AddWorkspace(configFilePath string, ws common.Workspace) error
	GetWorkspaceLocations() common.WorkspaceLocations
	GetWorkspaceLocation(name string) (common.WorkspaceLocation, error)
	SetProjectPath(wsName string, name string, into string) error
	ProjectExists(wsName string, name string) (bool, error)
}

type git interface {
	GetRepositoryRootFromPath(path string) (string, error)
	Clone(repoUrl string, target string) error
}

type Manager struct {
	fs            afero.Fs
	tui           tui
	configManager config
	git           git
}

func NewManager(fs afero.Fs, tui tui, configManager config, git git) *Manager {
	return &Manager{
		fs:            fs,
		tui:           tui,
		configManager: configManager,
		git:           git,
	}
}
