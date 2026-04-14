package tui

import (
	tea "charm.land/bubbletea/v2"
)

type results_screen struct{}

func CreateResultsScreen() CurrentScreen {
	return &results_screen{}
}

func (m *results_screen) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	return cmd
}

func (m *results_screen) View(w, h int) tea.View {
	v := tea.NewView("works results screen")

	v.AltScreen = true
	return v
}
