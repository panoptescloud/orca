package docker

import (
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/panoptescloud/orca/internal/common"
	"go.yaml.in/yaml/v4"
)

// TODO: guard against nil arguments
func (c *Compose) LoadSpec(ws *common.Workspace, p *common.Project) (*types.Project, error) {
	rawSpec, err := c.LoadRawSpec(ws, p)

	if err != nil {
		return nil, err
	}

	spec := &types.Project{}

	err = yaml.Unmarshal([]byte(rawSpec), spec)

	return spec, err
}
