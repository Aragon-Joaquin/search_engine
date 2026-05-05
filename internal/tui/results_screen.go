package tui

import (
	"fmt"
	"strconv"

	"search_engine/internal/blobs"
	"search_engine/tools"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/viewport"
)

const (
	MIN_THRESHOLD = 5
)

type results_screen struct {
	queryMade string
	bState    *itemsState

	// viewport render
	ready bool

	searchInput textinput.Model
	viewport    viewport.Model
	spinner     spinner.Model
}

type itemsState struct {
	isLoading bool
	items     []*blobs.Blob
}

func CreateResultsScreen(search_query string) (CurrentScreen, tea.Cmd) {
	itemsState := &itemsState{
		isLoading: true,
		items:     []*blobs.Blob{},
	}

	go func() {
		res := rep.UserMakeQuery(search_query)

		if len(res) == 0 {
			res = tools.CrawlIntoIndexer(search_query)
		}

		blobs := []*blobs.Blob{}

		for _, b := range res {
			if b.Score < float64(MIN_THRESHOLD)/100 {
				continue
			}
			blobs = append(blobs, b)
		}

		itemsState.isLoading = false
		itemsState.items = blobs
	}()

	ti := textinput.New()
	ti.Placeholder = "Search again!"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 40
	ti.SetWidth(120)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &results_screen{
		queryMade: search_query,
		bState:    itemsState,

		searchInput: ti,
		viewport:    viewport.New(),

		spinner: s,
	}, s.Tick
}

var (
	CURRENT_SELECTOR        int = 0
	HEADER_FOCUSEABLE_ITEMS int = 1
)

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

			if CURRENT_SELECTOR != 0 && len(m.bState.items) >= CURRENT_SELECTOR-1 {
				item := m.bState.items[CURRENT_SELECTOR-1]
				if item != nil {
					return changeCurrentScreen(CreateBlobInfoScreen(item, m.queryMade))
				}

			}

		case "down":
			if CURRENT_SELECTOR+1 < len(m.bState.items)+HEADER_FOCUSEABLE_ITEMS {
				CURRENT_SELECTOR = CURRENT_SELECTOR + 1
			}

		case "up":
			if CURRENT_SELECTOR-1 >= 0 {
				CURRENT_SELECTOR = CURRENT_SELECTOR - 1
			}

		case "q":
			if m.bState.isLoading {
				return changeCurrentScreen(CreateMainScreen())
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
			m.viewport.SetHeight(msg.Height - headerHeight - (showKeysLayout.GetHeight()))
		}

	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	if m.bState.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

var (
	// general
	list = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder(), true).
		Padding(0, 2).
		MaxWidth(150)

	// header
	headerTitle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			Bold(true)

	headerDate = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right).
			Foreground(lipgloss.BrightBlack)

	header = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("#444"))

	// info card
	infoUrl = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left).
		Foreground(lipgloss.Blue).
		Underline(true)

	infoScore = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right)

	// description
	bottomDescription = lipgloss.NewStyle().
				PaddingTop(1).
				Align(lipgloss.Left).
				MaxHeight(2)

	// extras
	noItems = lipgloss.NewStyle().Italic(true).Bold(true).MarginTop(4).Align(lipgloss.Center)
)

func (m *results_screen) View(w, h int) tea.View {
	// render blobs info
	itemsListed := []string{}

	fmt.Println(m.bState.items)
	if !m.bState.isLoading && len(m.bState.items) <= 0 {
		itemsListed = []string{
			noItems.Width(w).Render("No items available!"),
		}
	} else {
		itemsListed = m.bodyView(w, h)
	}

	// normal logic
	itemListedRendered := lipgloss.JoinVertical(
		lipgloss.Center,
		itemsListed...,
	)

	spin := ""
	if m.bState.isLoading {
		spin = fmt.Sprintf("\n\n%s Loading forever...press q to quit\n\n", m.spinner.View())
	}

	var v tea.View
	if !m.ready {
		v.SetContent("\n  Initializing...")
	} else {
		m.viewport.SetContent(itemListedRendered)
		v.SetContent(fmt.Sprintf("%s\n%s\n%s", spin, m.headerView(), m.viewport.View()))

	}

	return v
}

var titleSearch = lipgloss.NewStyle().Padding(0, 1).Margin(0, MARGIN_SIDES).Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.BrightBlack)

func (m *results_screen) headerView() string {
	if CURRENT_SELECTOR < HEADER_FOCUSEABLE_ITEMS {
		titleSearch = titleSearch.Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.Yellow)
		m.searchInput.Focus()
	} else {
		m.searchInput.Blur()
	}

	if !m.bState.isLoading {
		return lipgloss.JoinHorizontal(lipgloss.Center, titleSearch.Render(m.searchInput.View()))
	}

	return ""
}

func (m *results_screen) bodyView(w, h int) []string {
	itemsListed := []string{}
	for index, i := range m.bState.items {
		scoreParsed := int(i.Score * 100)

		listMargin := MARGIN_SIDES * 2
		list = list.
			BorderForeground(lipgloss.Red).
			Margin(1, listMargin).
			Width(w - (listMargin * 2))

		if CURRENT_SELECTOR >= HEADER_FOCUSEABLE_ITEMS-1 && index == (CURRENT_SELECTOR-HEADER_FOCUSEABLE_ITEMS) {
			list = list.Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.BrightYellow)
		}

		// NOTE: header card (title + date)
		formattedDate := i.Datetime.Format("2006/01/2")

		headerTitleStr := headerTitle.
			Width(list.GetWidth() - len(formattedDate) - (listMargin * 2)).
			Render(i.Title)

		headerDateStr := headerDate.Width(len(formattedDate)).Render(formattedDate)

		// NOTE: information card (url, + score)
		scoreToStr := strconv.Itoa(scoreParsed) + "% Match"

		infoUrlStr := infoUrl.
			Width(list.GetWidth() - len(scoreToStr) - listMargin*2).
			Hyperlink(i.URL).
			Render(i.URL)

		infoScoreStr := infoScore.
			Width(len(scoreToStr)).
			Foreground(AssignColorToScore(scoreParsed)).
			Render(scoreToStr)

		informationCard := lipgloss.NewStyle().Render(infoUrlStr, infoScoreStr)

		// NOTE: description (bottom)
		bottomDescriptionStr := bottomDescription.
			MaxWidth(list.GetWidth() - listMargin*2).
			Render(i.Description)

		// united all
		itemsListed = append(itemsListed,
			list.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					header.Render(headerTitleStr, headerDateStr),
					informationCard,
					bottomDescriptionStr,
				),
			),
		)
	}
	return itemsListed
}
