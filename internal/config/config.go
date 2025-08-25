// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

import "github.com/panoptescloud/orca/internal/common"

type ConfigLogging struct {
	Level  string
	Format string
}

type ConfigProject struct {
	Name string
	Path string
}

type ConfigWorkspace struct {
	Name     string
	Path     string
	Projects []ConfigProject
}

func (self *ConfigWorkspace) GetProject(name string) *ConfigProject {
	for _, project := range self.Projects {
		if project.Name == name {
			return &project
		}
	}

	return nil
}

func (self *ConfigWorkspace) GetProjectIndex(name string) int {
	for i, p := range self.Projects {
		if p.Name == name {
			return i
		}
	}

	return -1
}

func (self *ConfigWorkspace) ReplaceProject(project ConfigProject) {
	i := self.GetProjectIndex(project.Name)

	if i == -1 {
		self.Projects = append(self.Projects, project)
		return
	}

	self.Projects[i] = project
}

type Config struct {
	Logging          ConfigLogging
	Workspaces       []ConfigWorkspace
	CurrentWorkspace string `yaml:"currentWorkspace" mapstructure:"current_workspace"`
}

func (self *Config) WorkspaceExists(ws string) bool {
	for _, workspace := range self.Workspaces {
		if workspace.Name == ws {
			return true
		}
	}

	return false
}

func (self *Config) AddWorkspace(name string, path string) error {
	if self.WorkspaceExists(name) {
		return common.ErrWorkspaceAlreadyExists{
			Name: name,
		}
	}

	self.Workspaces = append(self.Workspaces, ConfigWorkspace{
		Name: name,
		Path: path,
	})

	return nil
}

func (self *Config) GetWorkspace(name string) (*ConfigWorkspace, error) {
	for _, ws := range self.Workspaces {
		if ws.Name == name {
			return &ws, nil
		}
	}

	return nil, common.ErrUnknownWorkspace{
		Workspace: name,
	}
}

func (self *Config) ProjectExists(wsName string, projectName string) (bool, error) {
	ws, err := self.GetWorkspace(wsName)

	if err != nil {
		return false, err
	}

	for _, p := range ws.Projects {
		if p.Name == projectName {
			return true, nil
		}
	}

	return false, nil
}

func (self *Config) GetWorkspaceIndex(name string) int {
	for i, ws := range self.Workspaces {
		if ws.Name == name {
			return i
		}
	}

	return -1
}

func (self *Config) ReplaceWorkspace(ws ConfigWorkspace) {
	i := self.GetWorkspaceIndex(ws.Name)

	if i == -1 {
		self.Workspaces = append(self.Workspaces, ws)
		return
	}

	self.Workspaces[i] = ws
}

func (self *Config) SetProjectPath(wsName string, name string, path string) error {
	ws, err := self.GetWorkspace(wsName)

	if err != nil {
		return err
	}

	project := ws.GetProject(name)

	if project == nil {
		project = &ConfigProject{
			Name: name,
		}
	}

	project.Path = path

	ws.ReplaceProject(*project)
	self.ReplaceWorkspace(*ws)

	return nil
}

func NewDefaultConfig() *Config {
	return &Config{
		Logging: ConfigLogging{
			Level:  "none",
			Format: "text",
		},
		Workspaces:       []ConfigWorkspace{},
		CurrentWorkspace: "",
	}
}
