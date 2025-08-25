package workspaces

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type InitialiseDTO struct {
	// Not yet implemented, but will allow you to supply a git url; then clone it
	SourceGitUrl      string
	SourceDirectory   string
	WorkspaceFileName string
	Into              string
}

func (self *Manager) loadConfigFromPath(path string) (*common.Workspace, error) {
	contents, err := afero.ReadFile(self.fs, path)

	if err != nil {
		return nil, err
	}

	cfg := &common.Workspace{}

	if err := yaml.Unmarshal(contents, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (self *Manager) Initialise(dto InitialiseDTO) error {
	workspaceConfigPath := fmt.Sprintf("%s/%s", dto.SourceDirectory, dto.WorkspaceFileName)
	exists, err := afero.Exists(self.fs, workspaceConfigPath)

	if err != nil {
		return self.tui.RecordIfError("Failed to find config, this is likely a bug!", err)
	}

	if !exists {
		return self.tui.RecordIfError(
			fmt.Sprintf("Could not find workspace config: %s", workspaceConfigPath),
			common.ErrFileNotFound{
				Path: workspaceConfigPath,
			},
		)
	}

	cfg, err := self.loadConfigFromPath(workspaceConfigPath)

	if err != nil {
		return self.tui.RecordIfError("Failed to open config file!", err)
	}

	err = self.configManager.AddWorkspace(workspaceConfigPath, *cfg)

	if err != nil {
		if _, ok := err.(common.ErrWorkspaceAlreadyExists); ok {
			self.tui.Error("Workspace already exists!")
			return err
		}

		return self.tui.RecordIfError("Failed to register workspace, this is likely a bug!", err)
	}

	self.tui.Success("Workspace initialised!")
	return nil
}
