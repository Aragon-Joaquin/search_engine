package tui

import (
	"fmt"
	"log"

	"search_engine/internal/repository"
	"search_engine/internal/utils"
	"search_engine/tools"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
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

func changeCurrentScreen(c CurrentScreen, cmd ...tea.Cmd) tea.Cmd {
	screen = c

	cmd = append(cmd, tea.RequestWindowSize)
	return tea.Batch(cmd...)
}

// needs to implement the tea.Model interface
type PTYModel struct {
	// internal state
	isSmall  bool
	quitting bool

	// client info
	term   string
	width  int
	height int

	// extras
	// time time.Time

	// ui info
	keys keyMap
	help help.Model
}

const (
	MIN_SCREEN_WIDTH  = 96
	MIN_SCREEN_HEIGHT = 30
)

// TODO: uhm...is there a better way?
var (
	rep *repository.Repository
	crw *tools.Crawler
)

func CreatePTYModel(r *repository.Repository, c *tools.Crawler, w, h int, t string) PTYModel {
	rep = r
	crw = c

	helpKeys := help.New()
	helpKeys.ShowAll = true

	pty := PTYModel{
		isSmall:  false,
		quitting: false,

		width:  w,
		height: h,
		term:   t,

		keys: initKeysMap,
		help: helpKeys,
	}

	return pty
}

func (m PTYModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
	)
}

var (
	min_size_screen_warn = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)
	min_text_title       = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Bold(true).
				Underline(true).
				Foreground(lipgloss.Yellow).
				Render("Screen is too small!")
	min_text_description = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Bold(true).
				Render(fmt.Sprintf("Minimum size: %dw x %dh", MIN_SCREEN_WIDTH, MIN_SCREEN_HEIGHT))

	min_paragraph = lipgloss.JoinVertical(lipgloss.Center,
		min_text_title,
		"\n",
		min_text_description,
	)
)

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
		if m.height <= MIN_SCREEN_HEIGHT && m.width <= MIN_SCREEN_WIDTH {
			m.isSmall = true
		} else {
			m.isSmall = false
		}
		log.Println("WIDTH: ", msg.Width, " HEIGHT: ", msg.Height)
		//TODO: screen gets overwritten on WindowSize event
		// i think the cause of this is the pattern im using

		return m, tea.Sequence(tea.ClearScreen, screen.Update(msg))

	default:
		log.Printf("UNKNOWN: %#v\n", msg)
	}

	return m, screen.Update(msg)
}

var showKeysLayout = lipgloss.NewStyle().Margin(0, 1).AlignVertical(lipgloss.Top).AlignHorizontal(lipgloss.Left).Height(4).Border(softBorder).BorderForeground(lipgloss.BrightBlack)

var softBorder = lipgloss.Border{
	Top:          "-",
	Bottom:       "-",
	Left:         "∣",
	Right:        "∣",
	TopLeft:      "┌",
	TopRight:     "┐",
	BottomLeft:   "└",
	BottomRight:  "┘",
	MiddleLeft:   "∣",
	MiddleRight:  "∣",
	Middle:       "∣",
	MiddleTop:    "-",
	MiddleBottom: "-",
}

func (m PTYModel) View() tea.View {
	if m.quitting {
		return tea.NewView("\n")
	}

	if m.isSmall {
		return tea.NewView(min_size_screen_warn.Width(m.width).Height(m.height).Render(min_paragraph))
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
