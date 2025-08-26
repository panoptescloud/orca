package workspaces

import (
	"fmt"
	"path/filepath"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/spf13/afero"
)

type CloneDTO struct {
	WorkspaceName string
	Project       string
	To            string
}

func (self *Manager) cloneAllProjects(ws *common.Workspace, into string) error {
	for _, project := range ws.Projects {
		dir := fmt.Sprintf("%s/%s", into, project.Name)

		if err := self.cloneSingleProject(ws, project, dir); err != nil {
			return err
		}
	}

	return nil
}

func (self *Manager) cloneSingleProject(ws *common.Workspace, project common.Project, into string) error {
	projectExists, err := self.configManager.ProjectExists(ws.Name, project.Name)

	if err != nil {
		return err
	}

	exists, err := afero.Exists(self.fs, into)

	if err != nil {
		return err
	}

	if exists {
		if projectExists {
			self.tui.Info(fmt.Sprintf("Skipping '%s' as it already exists!", project.Name))
			return nil
		}

		return common.ErrDirectoryAlreadyExists{
			Path: into,
		}
	}

	parentDir := filepath.Dir(into)

	if err := self.fs.MkdirAll(parentDir, 0755); err != nil {
		return err
	}

	self.tui.Info(fmt.Sprintf("Cloning '%s' into %s...", project.Repository.SSH, into))

	err = self.git.Clone(project.Repository.SSH, into)

	if err != nil {
		return self.tui.RecordIfError("Clone failed!", err)
	}

	return self.tui.RecordIfError(
		"Failed to save config, you will have to edit this manually but adding the path into your orca config!",
		self.configManager.SetProjectPath(ws.Name, project.Name, into),
	)
}

func (self *Manager) getCloneTargetDir(wsConfigPath string, target string) (string, error) {
	if target != "" {
		return filepath.Abs(target)
	}

	wsRoot, err := self.git.GetRepositoryRootFromPath(wsConfigPath)

	if err != nil {
		return "", err
	}

	return filepath.Abs(fmt.Sprintf("%s/..", wsRoot))
}

func (self *Manager) Clone(dto CloneDTO) error {

	location, err := self.configManager.GetWorkspaceLocation(dto.WorkspaceName)

	if err != nil {
		if _, ok := err.(common.ErrUnknownWorkspace); ok {
			return self.tui.RecordIfError(fmt.Sprintf("Unknown workspace: %s", dto.WorkspaceName), err)
		}

		return err
	}

	cfg, err := self.loadConfigFromPath(location.Path)

	if err != nil {
		return err
	}

	targetDir, err := self.getCloneTargetDir(location.Path, dto.To)

	if err != nil {
		return self.tui.RecordIfError("Failed to determine target directory!", err)
	}

	var cloneErr error = nil

	if dto.Project != "" {
		project, err := cfg.GetProject(dto.Project)

		if err != nil {
			return self.tui.RecordIfError("Failed to clone specified project!", err)
		}

		targetDir := targetDir

		if dto.To == "" {
			targetDir = fmt.Sprintf("%s/%s", targetDir, project.Name)
		}

		cloneErr = self.cloneSingleProject(cfg, *project, targetDir)
	} else {
		cloneErr = self.cloneAllProjects(cfg, targetDir)
	}

	if cloneErr != nil {
		if _, ok := cloneErr.(common.ErrDirectoryAlreadyExists); ok {
			return self.tui.RecordIfError("Cannot clone into the given directory, it already exists!", cloneErr)
		}

		return cloneErr
	}

	return nil
}
