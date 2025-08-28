package workspaces

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/spf13/afero"
)

type InitialiseDTO struct {
	// Not yet implemented, but will allow you to supply a git url; then clone it
	SourceGitUrl      string
	SourceDirectory   string
	WorkspaceFileName string
	Into              string
}

func (m *Manager) Initialise(dto InitialiseDTO) error {
	workspaceConfigPath := fmt.Sprintf("%s/%s", dto.SourceDirectory, dto.WorkspaceFileName)
	exists, err := afero.Exists(m.fs, workspaceConfigPath)

	if err != nil {
		return m.tui.RecordIfError("Failed to find config, this is likely a bug!", err)
	}

	if !exists {
		return m.tui.RecordIfError(
			fmt.Sprintf("Could not find workspace config: %s", workspaceConfigPath),
			common.ErrFileNotFound{
				Path: workspaceConfigPath,
			},
		)
	}

	cfg, err := m.workspaceRepo.LoadUnconfiguredWorkspace(workspaceConfigPath)

	if err != nil {
		return m.tui.RecordIfError("Failed to open config file!", err)
	}

	err = m.configManager.AddWorkspace(workspaceConfigPath, cfg.Name)

	if err != nil {
		if _, ok := err.(common.ErrWorkspaceAlreadyExists); ok {
			m.tui.Error("Workspace already exists!")
			return err
		}

		return m.tui.RecordIfError("Failed to register workspace, this is likely a bug!", err)
	}

	m.tui.Success("Workspace initialised!")
	return nil
}
