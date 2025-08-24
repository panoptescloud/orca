package config

import (
	"fmt"
	"path/filepath"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type tui interface {
	Error(msg ...string)
	Success(msg ...string)
}

type Manager struct {
	fs   afero.Fs
	path string
	tui  tui
	cfg  *Config
}

func (self *Manager) loadFromFile() (*Config, error) {
	contents, err := afero.ReadFile(self.fs, self.path)

	if err != nil {
		return nil, err
	}

	cfg := NewDefaultConfig()
	err = yaml.Unmarshal(contents, cfg)

	return cfg, err
}

func (self *Manager) SaveConfig(cfg *Config) error {
	contents, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	dir := filepath.Dir(self.path)

	if err := self.fs.MkdirAll(dir, 0755); err != nil {
		return err
	}

	err = afero.WriteFile(self.fs, self.path, contents, 0655)

	if err != nil {
		return err
	}

	self.cfg = cfg

	return nil
}

func (self *Manager) Get() *Config {
	return self.cfg
}

func (self *Manager) LoadOrCreate() (*Config, error) {
	exists, err := afero.Exists(self.fs, self.path)

	if err != nil {
		return nil, err
	}

	if !exists {
		cfg := NewDefaultConfig()
		self.cfg = cfg

		return cfg, self.SaveConfig(cfg)
	}

	cfg, err := self.loadFromFile()

	self.cfg = cfg

	return self.cfg, err
}

func (self *Manager) SwitchWorkspace(workspace string) error {
	cfg, err := self.loadFromFile()

	if err != nil {
		return err
	}

	if exists := cfg.WorkspaceExists(workspace); !exists {
		self.tui.Error(fmt.Sprintf("Workspace '%s' does not exist!", workspace))
		return common.ErrUnknownWorkspace{
			Workspace: workspace,
		}
	}

	cfg.CurrentWorkspace = workspace

	if err := self.SaveConfig(cfg); err != nil {
		return err
	}

	self.tui.Success(fmt.Sprintf("Switched to '%s' workspace", workspace))

	return nil
}

func (self *Manager) AddWorkspace(path string, ws *common.Workspace) error {
	cfg, err := self.loadFromFile()

	if err != nil {
		return err
	}

	err = cfg.AddWorkspace(ws.Name, path)

	if err != nil {
		return common.ErrWorkspaceAlreadyExists{
			Name: ws.Name,
		}
	}

	return self.SaveConfig(cfg)
}

func NewManager(tui tui, fs afero.Fs, path string) *Manager {
	return &Manager{
		fs:   fs,
		path: path,
		tui:  tui,
		cfg:  NewDefaultConfig(),
	}
}
