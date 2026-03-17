package tui

import (
	"fmt"
	"log"
	"search_engine/utils"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	APP_NAME_BANNER = `
     _                           
    / \   _ __ __ _  ___  _ __   
   / _ \ | '__/ _' |/ _ \| '_ \  
  / ___ \| | | (_| | (_) | | | | 
 /_/   \_\_|  \__, |\___/|_| |_| 
              |___/              `

	VERSION = utils.GetEnv(utils.ENV_VER)
)

// needs to implement the tea.Model interface
type PTYModel struct {
	// internal state
	// quitting bool

	// client info
	term   string
	width  int
	height int

	// extras
	// time time.Time

	// ui info
	spinner   spinner.Model
	textInput textinput.Model
}

func CreatePTYModel(w, h int, t string) PTYModel {
	ti := textinput.New()
	ti.Placeholder = "Search anything"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 40
	ti.SetWidth(120)

	return PTYModel{
		width:  w,
		height: h,
		term:   t,

		spinner:   spinner.New(),
		textInput: ti,
	}
}

func (m PTYModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
	)
}

func (m PTYModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			tea.ClearScreen()
			return m, tea.Quit
		}

		log.Printf("KEY PRESSED: %s\n", msg.Text)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		log.Println("WIDTH: ", msg.Width, " HEIGHT: ", msg.Height)
	default:
		log.Printf("UNKNOWN: %#v\n", msg)

	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
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

func (m PTYModel) View() tea.View {
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
	startX := (m.width - bodyWidth) / 2
	startY := (m.height - bodyHeight) / 2

	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y = (m.height / 2) - (textInputWrap.GetHeight() * 2)

		c.X += startX + 2                                 // 1 for padding + 1 for letter offset
		c.Y = startY + lipgloss.Height(renderedTitle) + 2 // 1 for padding + 1 margin top
	}

	var headerS strings.Builder
	headerS.WriteString(
		lipgloss.Place(m.width, 0, lipgloss.Center, lipgloss.Center,
			headerInfo.Render(
				fmt.Sprintf("v%s", VERSION),
			),
		),
	)

	v := tea.NewView(
		headerS.String() +
			"\n" +
			lipgloss.Place(
				m.width,
				m.height-(lipgloss.Height(headerS.String())*2),
				lipgloss.Center,
				lipgloss.Center,
				bodyS.String(),
			),
	)

	v.AltScreen = true
	v.Cursor = c
	return v
}
