package common

import "github.com/panoptescloud/orca/pkg/slices"

type LoaderPropertyCondition struct {
	Name  string
	Value any
}

type LoaderCondition struct {
	OS       string
	Arch     string
	Property LoaderPropertyCondition
}

type ExtraComposeFile struct {
	Path string
	When []LoaderCondition
}

type ComposeFiles struct {
	Primary string
	Extras  []ExtraComposeFile
}

type Property struct {
	Name    string
	Type    string
	Default any
}

type Extension struct {
	Name    string
	Chdir   string
	Command string
	Service string
}

type ProjectConfig struct {
	ComposeFiles ComposeFiles `yaml:"composeFiles"`
	Properties   []Property
	Hosts        []string
	Extensions   []Extension
}

type ProjectRepositoryConfig struct {
	SSH  string
	Self bool
}

type Project struct {
	Name             string
	RepositoryConfig ProjectRepositoryConfig
	ProjectDir       string
	Requires         []string
	Config           ProjectConfig
	IsRegistered     bool
}

func (p Project) GetName() string {
	return p.Name
}

func (p Project) GetKey() string {
	return p.Name
}

func (p Project) GetParents() []string {
	return p.Requires
}

func (p Project) GetChildren() []string {
	return []string{}
}

type Workspace struct {
	Name       string
	ConfigPath string
	Projects   []Project
}

func (ws *Workspace) GetProject(name string) (*Project, error) {
	p := slices.GetNamedElement(ws.Projects, name)

	if p == nil {
		return nil, ErrUnknownProject{
			Name: name,
		}
	}

	return p, nil
}

// --- unconfigured version

type UnconfiguredProject struct {
	Name             string
	RepositoryConfig ProjectRepositoryConfig
}

type UnconfiguredWorkspace struct {
	Name       string
	ConfigPath string
	Projects   []UnconfiguredProject
}
