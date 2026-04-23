package tui

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Up   key.Binding
	Down key.Binding

	Quit   key.Binding
	Select key.Binding

	Tab key.Binding
}

// both methods required for the interface
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Quit},
		{k.Tab},
	}
}

var initKeysMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("↹", "change focus"),
	),
}

