package workspaces

type ClearCurrentDTO struct {
}

func (m *Manager) ClearCurrent(dto ClearCurrentDTO) error {
	current := m.configManager.GetCurrentWorkspace()

	if current == "" {
		m.tui.Info("<no workspace selected>")
		return nil
	}

	return m.configManager.ClearCurrentWorkspace()
}
