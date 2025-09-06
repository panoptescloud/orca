package docker

import (
	"fmt"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/panoptescloud/orca/internal/common"
	"github.com/spf13/afero"
)

type overlayComposeParser interface {
	Parse(paths []string) (*types.Project, error)
}

type ComposeOverlayGenerator struct {
	fs         afero.Fs
	parser     overlayComposeParser
	overlayDir string
}

func emptyOverlay() *types.Project {
	return &types.Project{}
}
func (cog *ComposeOverlayGenerator) getOverlayPath(ws *common.Workspace, p *common.Project) string {
	return fmt.Sprintf("%s/%s/%s.yaml", cog.overlayDir, ws.Name, p.Name)
}

func (cog *ComposeOverlayGenerator) addNetwork(new *types.Project, original *types.Project, create bool) error {
	if new.Networks == nil {
		new.Networks = types.Networks{}
	}

	if create {
		new.Networks["orca"] = types.NetworkConfig{
			Name: "orca-ws",
		}
	} else {
		new.Networks["orca"] = types.NetworkConfig{
			Name:     "orca-ws",
			External: true,
		}
	}

	if len(original.Services) == 0 {
		// Nothing to do if there are no services
		return nil
	}

	if new.Services == nil {
		new.Services = types.Services{}
	}

	for name, _ := range original.Services {
		if _, ok := new.Services[name]; !ok {
			new.Services[name] = types.ServiceConfig{}
		}

		newSvc, _ := new.Services[name]

		if new.Services[name].Networks == nil {
			newSvc.Networks = map[string]*types.ServiceNetworkConfig{}
		}

		newSvc.Networks["orca"] = &types.ServiceNetworkConfig{}

		new.Services[name] = newSvc
	}

	return nil
}

func (cog *ComposeOverlayGenerator) CreateOrRetrieve(ws *common.Workspace, p *common.Project) (string, error) {
	composeProject, err := cog.parser.Parse([]string{p.Config.ComposeFiles.Primary})
	if err != nil {
		return "", err
	}

	// fmt.Printf("ORIGINAL\n\n%#v\n\n", composeProject)

	o := emptyOverlay()

	if ws.OverlayConfig.Network.Enabled {
		if err := cog.addNetwork(o, composeProject, p.Name == ws.OverlayConfig.Network.CreateIn); err != nil {
			return "", err
		}
	}

	newFile, err := o.MarshalYAML()

	if err != nil {
		return "", err
	}

	overlayPath := cog.getOverlayPath(ws, p)
	overlayDir := filepath.Dir(overlayPath)

	if err := cog.fs.MkdirAll(overlayDir, 0700); err != nil {
		return "", err
	}

	if err := afero.WriteFile(cog.fs, overlayPath, newFile, 0655); err != nil {
		return "", err
	}

	// fmt.Printf("NEW\n\n%s\n\n", string(newFile))
	// if ws.

	// fmt.Printf("overlay:\n%#v\n\n", ws.OverlayConfig)

	return overlayPath, nil
}

func NewComposeOverlayGenerator(fs afero.Fs, parser overlayComposeParser, overlayDir string) *ComposeOverlayGenerator {
	return &ComposeOverlayGenerator{
		fs:         fs,
		parser:     parser,
		overlayDir: overlayDir,
	}
}
