package workspaces

import "fmt"

type LsDTO struct {
}

func (self *Manager) Ls(dto LsDTO) error {
	locs := self.configManager.GetWorkspaceLocations()

	if len(locs) == 0 {
		self.tui.Error("No workspaces exist!")

		return nil
	}

	for _, loc := range locs {
		self.tui.Info(fmt.Sprintf("%s -> %s", loc.Name, loc.Path))
	}

	return nil
}
