package config

import (
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Manager struct {
	fs   afero.Fs
	path string
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

	return afero.WriteFile(self.fs, self.path, contents, 0655)
}

func (self *Manager) LoadOrCreate() (*Config, error) {
	exists, err := afero.Exists(self.fs, self.path)

	if err != nil {
		return nil, err
	}

	if !exists {
		cfg := NewDefaultConfig()
		return cfg, self.SaveConfig(cfg)
	}

	return self.loadFromFile()
}

func NewManager(fs afero.Fs, path string) *Manager {
	return &Manager{
		fs:   fs,
		path: path,
	}
}
