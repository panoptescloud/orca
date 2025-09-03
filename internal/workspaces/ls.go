package workspaces

import "fmt"

type LsDTO struct {
}

func (m *Manager) Ls(dto LsDTO) error {
	locs := m.configManager.GetAllWorkspaceMeta()

	if len(locs) == 0 {
		m.tui.Error("No workspaces exist!")

		return nil
	}

	for _, loc := range locs {
		m.tui.Info(fmt.Sprintf("%s -> %s", loc.Name, loc.Path))
	}

	return nil
}
