// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

import "github.com/panoptescloud/orca/internal/common"

type ConfigLogging struct {
	Level  string
	Format string
}

type ConfigWorkspace struct {
	Name string
	Path string
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
