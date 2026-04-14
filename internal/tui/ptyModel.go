package tui

import (
	"log"

	"search_engine/internal/utils"

	"charm.land/bubbles/v2/spinner"
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

	VERSION = utils.GetEnv(utils.ENV_VER)
)

type CurrentScreen interface {
	Update(msg tea.Msg) tea.Cmd
	View(w, h int) tea.View
}

// i tried oop but the PTYModels keeps recreating so i loose the reference
var screen CurrentScreen = CreateMainScreen()

func changeCurrentScreen(c CurrentScreen) tea.Cmd {
	screen = c
	return tea.ClearScreen
}

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
	spinner spinner.Model
}

func CreatePTYModel(w, h int, t string) PTYModel {
	pty := PTYModel{
		width:  w,
		height: h,
		term:   t,

		spinner: spinner.New(),
	}

	return pty
}

func (m PTYModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
	)
}

func (m PTYModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			tea.ClearScreen()
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

func (m PTYModel) View() tea.View {
	return screen.View(m.width, m.height)
}
