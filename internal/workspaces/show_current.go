package workspaces

type ShowCurrentDTO struct {
}

func (m *Manager) ShowCurrent(dto ShowCurrentDTO) error {
	current := m.configManager.GetCurrentWorkspace()

	if current == "" {
		m.tui.Info("<no workspace selected>")
		return nil
	}

	m.tui.Info(current)
	return nil
}
