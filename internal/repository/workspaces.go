package repository

import (
	"fmt"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/repository/internal/model"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type userConfig interface {
	GetWorkspaceMeta(name string) (common.WorkspaceMeta, error)
	GetProjectMeta(wsName string, projectName string) (common.ProjectMeta, error)
}

func convertLoaderCondition(e model.LoaderCondition) common.LoaderCondition {
	return common.LoaderCondition{
		OS:   e.OS,
		Arch: e.Arch,
		Property: common.LoaderPropertyCondition{
			Name:  e.Property.Name,
			Value: e.Property.Value,
		},
	}
}

func convertLoaderConditions(cfgConds []model.LoaderCondition) []common.LoaderCondition {
	conds := make([]common.LoaderCondition, len(cfgConds))
	for i, c := range cfgConds {
		conds[i] = convertLoaderCondition(c)
	}

	return conds
}

func convertExtraComposeFile(e model.ExtraComposeFile) common.ExtraComposeFile {
	return common.ExtraComposeFile{
		Path: e.Path,
		When: convertLoaderConditions(e.When),
	}
}

func convertExtraComposeFiles(files []model.ExtraComposeFile) []common.ExtraComposeFile {
	extras := make([]common.ExtraComposeFile, len(files))

	for i, e := range files {
		extras[i] = convertExtraComposeFile(e)
	}

	return extras
}

func convertComposeFiles(cfg model.ComposeFiles) common.ComposeFiles {
	return common.ComposeFiles{
		Primary: cfg.Primary,
		Extras:  convertExtraComposeFiles(cfg.Extras),
	}
}

func convertProperty(p model.Property) common.Property {
	return common.Property{
		Name:    p.Name,
		Type:    p.Type,
		Default: p.Default,
	}
}

func convertProperties(cfgProps []model.Property) []common.Property {
	props := make([]common.Property, len(cfgProps))

	for i, p := range cfgProps {
		props[i] = convertProperty(p)
	}

	return props
}

func convertExtension(e model.Extension) common.Extension {
	return common.Extension{
		Name:    e.Name,
		Chdir:   e.Chdir,
		Command: e.Command,
		Service: e.Service,
	}
}

func convertExtensions(cfgExts []model.Extension) []common.Extension {
	exts := make([]common.Extension, len(cfgExts))

	for i, e := range cfgExts {
		exts[i] = convertExtension(e)
	}

	return exts
}

func buildProject(wsPCfg model.WorkspaceProjectConfig, pCfg model.ProjectConfig) common.Project {
	return common.Project{
		Name: wsPCfg.Name,
		RepositoryConfig: common.ProjectRepositoryConfig{
			SSH:  wsPCfg.Repository.SSH,
			Self: wsPCfg.Repository.Self,
		},
		IsRegistered: true,
		ProjectDir:   wsPCfg.Path,
		Requires:     wsPCfg.Requires,
		Config: common.ProjectConfig{
			ComposeFiles:    convertComposeFiles(pCfg.ComposeFiles),
			Properties:      convertProperties(pCfg.Properties),
			Hosts:           pCfg.Hosts,
			TLSCertificates: pCfg.TLSCertificates,
			Extensions:      convertExtensions(pCfg.Extensions),
		},
	}
}

type WorkspaceRepository struct {
	fs         afero.Fs
	userConfig userConfig
}

func (wcr *WorkspaceRepository) loadProject(projectMeta common.ProjectMeta) (*model.ProjectConfig, error) {
	path := fmt.Sprintf("%s/%s", strings.TrimSuffix(projectMeta.Path, "/"), common.DefaultProjectFileName)

	found, err := afero.Exists(wcr.fs, path)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, common.ErrProjectConfigNotFound{
			Path: path,
		}
	}

	contents, err := afero.ReadFile(wcr.fs, path)

	if err != nil {
		return nil, err
	}

	cfg := &model.ProjectConfig{}

	if err := yaml.Unmarshal(contents, cfg); err != nil {
		return nil, common.ErrConfigIsNotValidYAML{
			Path:       path,
			ParseError: err.Error(),
		}
	}

	return cfg, nil
}

func (wcr *WorkspaceRepository) loadConfigFromPath(path string) (*model.WorkspaceConfig, error) {
	found, err := afero.Exists(wcr.fs, path)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, common.ErrWorkspaceConfigNotFound{
			Path: path,
		}
	}

	contents, err := afero.ReadFile(wcr.fs, path)

	if err != nil {
		return nil, err
	}

	cfg := &model.WorkspaceConfig{}

	if err := yaml.Unmarshal(contents, cfg); err != nil {
		return nil, common.ErrConfigIsNotValidYAML{
			Path:       path,
			ParseError: err.Error(),
		}
	}

	return cfg, nil
}

func (wcr *WorkspaceRepository) LoadUnconfiguredWorkspace(path string) (*common.UnconfiguredWorkspace, error) {
	ws, err := wcr.loadConfigFromPath(path)

	if err != nil {
		return nil, err
	}

	projects := make([]common.UnconfiguredProject, len(ws.Projects))

	for i, p := range ws.Projects {
		projects[i] = common.UnconfiguredProject{
			Name: p.Name,
			RepositoryConfig: common.ProjectRepositoryConfig{
				SSH:  p.Repository.SSH,
				Self: p.Repository.Self,
			},
		}
	}

	return &common.UnconfiguredWorkspace{
		Name:       ws.Name,
		ConfigPath: path,
		Projects:   projects,
	}, nil
}

func (wcr *WorkspaceRepository) Load(name string) (*common.Workspace, error) {
	wsMeta, err := wcr.userConfig.GetWorkspaceMeta(name)
	if err != nil {
		return nil, err
	}

	cfg, err := wcr.loadConfigFromPath(wsMeta.Path)

	if err != nil {
		return nil, err
	}

	ws := &common.Workspace{
		Name:       cfg.Name,
		ConfigPath: wsMeta.Path,
		Projects:   make([]common.Project, len(cfg.Projects)),
	}

	for i, pCfg := range cfg.Projects {
		projectMeta, err := wcr.userConfig.GetProjectMeta(wsMeta.Name, pCfg.Name)

		if err != nil {
			if _, ok := err.(common.ErrUnknownProject); ok {
				ws.Projects[i] = common.Project{
					Name: pCfg.Name,
					RepositoryConfig: common.ProjectRepositoryConfig{
						SSH:  pCfg.Repository.SSH,
						Self: pCfg.Repository.Self,
					},
					IsRegistered: false,
					ProjectDir:   "",
					Requires:     pCfg.Requires,
					Config:       common.ProjectConfig{},
				}
				continue
			}

			return nil, err
		}

		p, err := wcr.loadProject(projectMeta)

		if err != nil {
			return nil, err
		}

		ws.Projects[i] = buildProject(pCfg, *p)
	}

	return ws, nil
}

func NewWorkspaceRepository(fs afero.Fs, userConfig userConfig) *WorkspaceRepository {
	return &WorkspaceRepository{
		fs:         fs,
		userConfig: userConfig,
	}
}
