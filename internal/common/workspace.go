package common

import (
	stdslices "slices"

	"github.com/panoptescloud/orca/pkg/slices"
)

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
	ComposeFiles    ComposeFiles
	Properties      []Property
	Hosts           []string
	TLSCertificates []string
	Extensions      []Extension
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

func (ws *Workspace) GetUniqueTLSCertificates() []string {
	certs := []string{}

	for _, p := range ws.Projects {
		for _, c := range p.Config.TLSCertificates {
			if !stdslices.Contains(certs, c) {
				certs = append(certs, c)
			}
		}
	}

	return certs
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
