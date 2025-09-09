package controller

import "github.com/panoptescloud/orca/internal/common"

type ShowComposeCommandDTO struct {
	Workspace string
	Project   string
}

func (c *Controller) ShowComposeCommand(dto ShowComposeCommandDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	if ctx.Project == nil {
		return common.ErrInvalidExecutionContext{
			Msg: "this command must be run in a singular project context",
		}
	}

	return c.compose.ShowCommand(ctx.Workspace, ctx.Project)
}
