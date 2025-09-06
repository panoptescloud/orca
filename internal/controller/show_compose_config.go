package controller

import "github.com/panoptescloud/orca/internal/common"

type ShowComposeConfigDTO struct {
	Workspace string
	Project   string
}

func (c *Controller) ShowComposeConfig(dto ShowComposeConfigDTO) error {
	ctx, err := c.resolveContext(dto.Workspace, dto.Project)

	if err != nil {
		return err
	}

	if ctx.Project == nil {
		return common.ErrInvalidExecutionContext{
			Msg: "this command must be run in a singular project context",
		}
	}

	return c.compose.Show(ctx.Workspace, ctx.Project)
}
