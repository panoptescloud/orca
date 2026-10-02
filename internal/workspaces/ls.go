package workspaces

type LsDTO struct {
}

func (m *Manager) Ls(dto LsDTO) error {
	current := m.configManager.GetCurrentWorkspace()
	locs := m.configManager.GetAllWorkspaceMeta()

	if len(locs) == 0 {
		m.tui.Error("No workspaces exist!")

		return nil
	}

	rows := make([][]string, len(locs))

	for i, loc := range locs {
		inUse := "No"

		if loc.Name == current {
			inUse = "Yes"
		}
		rows[i] = []string{
			inUse,
			loc.Name,
			loc.Path,
		}
	}

	m.tui.Table([]string{
		"In Use?",
		"Project",
		"Path",
	}, rows)

	return nil
}
