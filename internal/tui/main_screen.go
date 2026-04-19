package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type main_screen struct {
	textInput textinput.Model
}

func CreateMainScreen() CurrentScreen {
	ti := textinput.New()
	ti.Placeholder = "Search anything"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 40
	ti.SetWidth(120)

	return &main_screen{
		textInput: ti,
	}
}

func (m *main_screen) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			value := strings.TrimSpace(m.textInput.Value())
			if m.textInput.Value() != "" {
				return changeCurrentScreen(CreateResultsScreen(value))
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

var (
	titleModel = lipgloss.NewStyle().
			Align(lipgloss.Center).
			BorderBottom(true).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Red)

	textInputWrap = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.BrightBlack).
			Padding(0, 1).
			Width(50).
			MarginTop(1)

	headerInfo = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(70).
			Margin(2, 10).
			Height(1).
			BorderStyle(lipgloss.ASCIIBorder()).
			BorderForeground(lipgloss.Red)
)

func (m *main_screen) View(w, h int) tea.View {
	var bodyS strings.Builder
	var c *tea.Cursor

	renderedTitle := titleModel.Render(APP_NAME_BANNER)
	bodyS.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Center,
			// title of the app
			renderedTitle,

			// text input
			textInputWrap.Render(
				m.textInput.View(),
			),
		))

	// get window size
	bodyWidth, bodyHeight := lipgloss.Size(bodyS.String())

	// get starting window points
	startX := (w - bodyWidth) / 2
	startY := (h - bodyHeight) / 2

	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y = (h / 2) - (textInputWrap.GetHeight() * 2)

		c.X += startX + 2                                 // 1 for padding + 1 for letter offset
		c.Y = startY + lipgloss.Height(renderedTitle) + 2 // 1 for padding + 1 margin top
	}

	var headerS strings.Builder
	headerS.WriteString(
		lipgloss.Place(w, 0, lipgloss.Center, lipgloss.Center,
			headerInfo.Render(
				fmt.Sprintf("v%s", VERSION),
			),
		),
	)

	v := tea.NewView(
		headerS.String() +
			"\n" +
			lipgloss.Place(
				w,
				h-(lipgloss.Height(headerS.String())*2),
				lipgloss.Center,
				lipgloss.Center,
				bodyS.String(),
			),
	)

	v.AltScreen = true
	v.Cursor = c
	return v
}
