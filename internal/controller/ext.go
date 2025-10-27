package controller

import (
	"errors"
	"strings"

	"github.com/panoptescloud/orca/internal/common"
)

type ExecuteExtensionDTO struct {
	Workspace string
	Project   string
	Name      string
	Args      []string
}

// TODO: de-dupe this by reusing the ExecOrRun method
func (c *Controller) executeExtensionInService(dto ExecuteExtensionDTO, ctx runtimeContext, ext common.Extension) error {
	cmdArgs := strings.Split(ext.Command, " ")

	if len(dto.Args) > 0 {
		cmdArgs = append(cmdArgs, dto.Args...)
	} else {
		cmdArgs = append(cmdArgs, ext.DefaultArgs...)
	}

	isRunning, err := c.compose.IsSvcRunning(ctx.Workspace, ctx.Project, ext.Service)

	if err != nil {
		return err
	}

	if !isRunning {
		return c.compose.Run(ctx.Workspace, ctx.Project, ext.Service, cmdArgs)
	}

	return c.compose.Exec(ctx.Workspace, ctx.Project, ext.Service, cmdArgs)
}

func (c *Controller) ExecuteExtension(dto ExecuteExtensionDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	if ctx.Project == nil {
		return common.ErrInvalidExecutionContext{
			Msg: "Extensions must be executed within a project context!",
		}
	}

	ext, err := ctx.Project.FindExtension(dto.Name)

	if err != nil {
		return c.tui.RecordIfError("Extension does not exist in project context.", err)
	}

	if ext.Service != "" {
		return c.executeExtensionInService(dto, ctx, ext)
	}

	return c.tui.RecordIfError("Extension must have a service defined", errors.New("alternative is not yet implemented"))
}
