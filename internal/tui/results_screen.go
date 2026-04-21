package tui

import (
	"fmt"
	"image/color"
	"strconv"

	"search_engine/internal/blobs"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/viewport"
)

const (
	MIN_THRESHOLD = -1
)

type results_screen struct {
	items []*blobs.Blob

	ready bool

	searchInput textinput.Model
	viewport    viewport.Model
}

func CreateResultsScreen(search_query string) CurrentScreen {
	ti := textinput.New()
	ti.Placeholder = "Search again!"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 40
	ti.SetWidth(120)

	minBlobs := []*blobs.Blob{}
	for _, b := range rep.UserMakeQuery(search_query) {
		if b.Score < float64(MIN_THRESHOLD)/100 {
			continue
		}
		minBlobs = append(minBlobs, b)
	}
	return &results_screen{
		items:       minBlobs,
		searchInput: ti,
		viewport:    viewport.New(),
	}
}

func AssignColorToScore(s int) color.Color {
	if s < 10 {
		return lipgloss.BrightBlack
	}

	if s >= 10 && s < 20 {
		return lipgloss.Red
	}

	if s >= 20 && s < 30 {
		return lipgloss.Yellow
	}

	if s >= 30 && s < 40 {
		return lipgloss.Green
	}

	if s >= 40 && s < 50 {
		return lipgloss.Cyan
	}

	return lipgloss.Magenta
}

var CURRENT_SELECTOR int = 0

func (m *results_screen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if ok := m.searchInput.Focused(); ok {
				if len(m.searchInput.Value()) > 3 {
					return changeCurrentScreen(CreateResultsScreen(m.searchInput.Value()))
				}
			}

		case "down", "j":
			if CURRENT_SELECTOR+1 < len(m.items) {
				CURRENT_SELECTOR = CURRENT_SELECTOR + 1
			}

		case "up", "k":
			if CURRENT_SELECTOR-1 >= 0 {
				CURRENT_SELECTOR = CURRENT_SELECTOR - 1
			}

		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-headerHeight))
			m.viewport.YPosition = headerHeight

			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - headerHeight)
		}

	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *results_screen) View(w, h int) tea.View {
	itemsListed := []string{}
	// render blobs info
	for index, i := range m.items {
		scoreParsed := int(i.Score * 100)

		listMargin := MARGIN_SIDES * 2
		list := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Margin(1, listMargin).
			Width(w-(listMargin*2)).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Red).
			Padding(0, 2).
			MaxWidth(150)

		descriptionFallback := "No description provided"
		if i.Description != "" {
			descriptionFallback = i.Description
		}

		if index == CURRENT_SELECTOR {
			list = list.Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.BrightYellow)
		}

		// NOTE: header card (title + date)
		formattedDate := i.Datetime.Format("2006/01/2")

		headerTitle := lipgloss.NewStyle().
			Width(list.GetWidth() - len(formattedDate) - (listMargin * 2)).
			AlignHorizontal(lipgloss.Left).
			Bold(true).
			Render(i.Title)

		headerDate := lipgloss.NewStyle().
			Width(len(formattedDate)).
			AlignHorizontal(lipgloss.Right).
			Foreground(lipgloss.BrightBlack).
			Render(formattedDate)

		header := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("#444")).Render(headerTitle, headerDate)

		// NOTE: information card (url, + score)
		scoreToStr := strconv.Itoa(scoreParsed) + "% Match"

		infoUrl := lipgloss.NewStyle().
			Width(list.GetWidth() - len(scoreToStr) - listMargin*2).
			AlignHorizontal(lipgloss.Left).
			Hyperlink(i.URL).
			Foreground(lipgloss.Blue).
			Underline(true).
			Render(i.URL)

		infoScore := lipgloss.NewStyle().
			Width(len(scoreToStr)).
			AlignHorizontal(lipgloss.Right).
			Foreground(AssignColorToScore(scoreParsed)).
			Render(scoreToStr)

		informationCard := lipgloss.NewStyle().Render(infoUrl, infoScore)

		// description (bottom)
		bottomDescription := lipgloss.NewStyle().
			PaddingTop(1).
			Align(lipgloss.Left).
			Render(descriptionFallback)

		// united all
		itemsListed = append(itemsListed,
			list.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					header,
					informationCard,
					bottomDescription,
				),
			),
		)

	}

	// normal logic
	itemListedRendered := lipgloss.JoinVertical(
		lipgloss.Center,
		itemsListed...,
	)

	var v tea.View
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	if !m.ready {
		v.SetContent("\n  Initializing...")
	} else {
		m.viewport.SetContent(itemListedRendered)
		v.SetContent(fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View()))
	}

	return v
}

func (m *results_screen) headerView() string {
	title := lipgloss.NewStyle().Margin(1, MARGIN_SIDES*2).Render(m.searchInput.View())
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}
