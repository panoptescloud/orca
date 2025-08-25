package config

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/pkg/slices"
)

type configLogging struct {
	Level  string
	Format string
}

type configProject struct {
	Name string
	Path string
}

func (self configProject) GetName() string {
	return self.Name
}

type configWorkspace struct {
	Name     string
	Path     string
	Projects []configProject
}

func (self configWorkspace) GetName() string {
	return self.Name
}

type config struct {
	Logging          configLogging
	Workspaces       []configWorkspace
	CurrentWorkspace string `yaml:"currentWorkspace" mapstructure:"current_workspace"`
}

func (self *config) workspaceExists(name string) bool {
	return slices.NamedElementExists(self.Workspaces, name)
}

func (self *config) getWorkspace(name string) (configWorkspace, error) {
	v := slices.GetNamedElement(self.Workspaces, name)

	if v == nil {
		return configWorkspace{}, common.ErrUnknownWorkspace{
			Name: name,
		}
	}

	return *v, nil
}

func (self *config) projectExists(wsName string, name string) (bool, error) {
	ws, err := self.getWorkspace(wsName)

	if err != nil {
		return false, err
	}

	return slices.NamedElementExists(ws.Projects, name), nil
}

func (self *config) getProject(wsName string, name string) (configProject, error) {
	ws, err := self.getWorkspace(wsName)

	if err != nil {
		return configProject{}, err
	}

	v := slices.GetNamedElement(ws.Projects, name)

	if v == nil {
		return configProject{}, common.ErrUnknownProject{
			Name: name,
		}
	}

	return *v, nil
}
