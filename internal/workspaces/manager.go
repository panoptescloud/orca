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
	AddWorkspace(configFilePath string, name string) error
	GetAllWorkspaceMeta() []common.WorkspaceMeta
	GetWorkspaceMeta(name string) (common.WorkspaceMeta, error)
	SetProjectPath(wsName string, name string, into string) error
	ProjectExists(wsName string, name string) (bool, error)
}

type git interface {
	GetRepositoryRootFromPath(path string) (string, error)
	Clone(repoUrl string, target string) error
}

type workspaceRepo interface {
	Load(name string) (*common.Workspace, error)
	LoadUnconfiguredWorkspace(path string) (*common.UnconfiguredWorkspace, error)
}

type Manager struct {
	fs            afero.Fs
	tui           tui
	configManager config
	git           git
	workspaceRepo workspaceRepo
}

// func (self *Manager) loadConfigFromPath(path string) (*common.WorkspaceMeta, error) {
// 	contents, err := afero.ReadFile(self.fs, path)

// 	if err != nil {
// 		return nil, err
// 	}

// 	cfg := &common.WorkspaceMeta{}

// 	if err := yaml.Unmarshal(contents, cfg); err != nil {
// 		return nil, err
// 	}

// 	return cfg, nil
// }

// func (self *Manager) LoadWorkspaceMeta(loc common.WorkspaceMeta) (*common.WorkspaceMeta, error) {
// 	return self.loadConfigFromPath(loc.Path)
// }

func NewManager(fs afero.Fs, tui tui, configManager config, git git, workspaceRepo workspaceRepo) *Manager {
	return &Manager{
		fs:            fs,
		tui:           tui,
		configManager: configManager,
		git:           git,
		workspaceRepo: workspaceRepo,
	}
}
