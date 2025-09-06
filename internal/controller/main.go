package controller

import (
	"os"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
)

type config interface {
	GetAllProjectMeta() []common.ProjectMeta
	GetCurrentWorkspace() string
	GetWorkspaceMeta(name string) (common.WorkspaceMeta, error)
}

type workspaceRepository interface {
	Load(path string) (*common.Workspace, error)
}

type compose interface {
	Up(*common.Workspace, *common.Project) error
	Down(ws *common.Workspace, p *common.Project) error
	Show(ws *common.Workspace, p *common.Project) error
	Exec(ws *common.Workspace, p *common.Project, service string, cmdArgs []string) error
}

type Controller struct {
	cfg           config
	workspaceRepo workspaceRepository
	compose       compose
}

type runtimeContext struct {
	Workspace *common.Workspace
	Project   *common.Project
}

func (c *Controller) getProjectFromWorkdir() (*common.ProjectMeta, error) {
	dir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	projects := c.cfg.GetAllProjectMeta()

	for _, p := range projects {
		if strings.HasPrefix(dir, p.Path) {
			return &p, nil
		}
	}

	return nil, nil
}

func (c *Controller) buildRuntimeContext(wsName string, projectName string) (runtimeContext, error) {
	meta, err := c.cfg.GetWorkspaceMeta(wsName)

	if err != nil {
		return runtimeContext{}, err
	}

	ws, err := c.workspaceRepo.Load(meta.Name)

	if err != nil {
		return runtimeContext{}, err
	}

	if projectName == "" {
		return runtimeContext{
			Workspace: ws,
		}, nil
	}

	project, err := ws.GetProject(projectName)

	if err != nil {
		return runtimeContext{}, err
	}

	return runtimeContext{
		Workspace: ws,
		Project:   project,
	}, nil
}

func (c *Controller) resolveContext(ws string, project string) (runtimeContext, error) {
	// We've got a specific workspace and project
	if ws != "" && project != "" {
		return c.buildRuntimeContext(ws, project)
	}

	// No workspace, but a project was specified, so we'll assume that it's the
	// current workspace
	if ws == "" && project != "" {
		return c.buildRuntimeContext(c.cfg.GetCurrentWorkspace(), project)
	}

	// The workspace was specified, but not project was so we'll use that workspace
	// and assume all projects.
	if ws != "" && project == "" {
		return c.buildRuntimeContext(ws, "")
	}

	// If we get here, both of the options were empty, so we'll try resolve from
	// the working directory
	p, err := c.getProjectFromWorkdir()

	if err != nil {
		return runtimeContext{}, err
	}

	if p != nil {
		return c.buildRuntimeContext(p.WorkspaceName, p.Name)
	}

	// We weren't in a project directory, so our final resort is to use the
	// current workspace, and assume all projects should be started.
	return c.buildRuntimeContext(c.cfg.GetCurrentWorkspace(), "")
}

func NewController(cfg config, wsRepo workspaceRepository, compose compose) *Controller {
	return &Controller{
		cfg:           cfg,
		workspaceRepo: wsRepo,
		compose:       compose,
	}
}
