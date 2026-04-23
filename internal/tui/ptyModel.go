package tui

import (
	"log"

	"search_engine/internal/repository"
	"search_engine/internal/utils"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"

	tea "charm.land/bubbletea/v2"
)

var (
	APP_NAME_BANNER = `
     _                           
    / \   _ __ __ _  ___  _ __   
   / _ \ | '__/ _' |/ _ \| '_ \  
  / ___ \| | | (_| | (_) | | | | 
 /_/   \_\_|  \__, |\___/|_| |_| 
              |___/              `

	VERSION      = utils.GetEnv(utils.ENV_VER)
	MARGIN_SIDES = 2
)

type CurrentScreen interface {
	Update(msg tea.Msg) tea.Cmd
	View(w, h int) tea.View
}

// i tried oop but the PTYModels keeps recreating so i loose the reference
var screen CurrentScreen = CreateMainScreen()

func changeCurrentScreen(c CurrentScreen) tea.Cmd {
	screen = c
	return tea.Batch(tea.RequestWindowSize)
}

// needs to implement the tea.Model interface
type PTYModel struct {
	// internal state
	quitting bool

	// client info
	term   string
	width  int
	height int

	// extras
	// time time.Time

	// ui info
	keys    keyMap
	spinner spinner.Model
	help    help.Model
}

// uhm...is there a better way?
var rep *repository.Repository

func CreatePTYModel(r *repository.Repository, w, h int, t string) PTYModel {
	rep = r

	helpKeys := help.New()
	helpKeys.ShowAll = true

	pty := PTYModel{
		quitting: false,

		width:  w,
		height: h,
		term:   t,

		keys:    initKeysMap,
		spinner: spinner.New(),
		help:    helpKeys,
	}

	return pty
}

func (m PTYModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
	)
}

func (m PTYModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// i dont know how to clear the screen on exit
	// without copying and pasting this everywhere
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
		log.Printf("KEY PRESSED: %s\n", msg.Text, msg.String())
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		log.Println("WIDTH: ", msg.Width, " HEIGHT: ", msg.Height)

	default:
		log.Printf("UNKNOWN: %#v\n", msg)
	}

	return m, screen.Update(msg)
}

var showKeysLayout = lipgloss.NewStyle().Margin(0, 2).AlignVertical(lipgloss.Top).AlignHorizontal(lipgloss.Left).Height(3)

func (m PTYModel) View() tea.View {
	if m.quitting {
		return tea.NewView("\n")
	}

	content := screen.View(m.width, m.height)
	content.AltScreen = true
	content.MouseMode = tea.MouseModeCellMotion
	content.WindowTitle = "Argon"

	layerMain := lipgloss.NewLayer(content.Content)
	layerKey := lipgloss.NewLayer(
		showKeysLayout.Render(m.help.View(m.keys)),
	).X(0).Y(m.height - (showKeysLayout.GetHeight()))

	content.SetContent(lipgloss.NewCompositor(layerMain, layerKey).Render())

	return content
}
