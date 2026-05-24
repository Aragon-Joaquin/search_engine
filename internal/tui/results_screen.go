package tui

import (
	"fmt"
	"strconv"
	"strings"

	"search_engine/internal/blobs"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/viewport"
)

type results_screen struct {
	queryMade string

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

var bState = &itemsState{
	isLoading: true,
	items:     []*blobs.Blob{},
}

func CreateResultsScreen(search_query string) (CurrentScreen, tea.Cmd) {
	bState.isLoading = true
	bState.items = []*blobs.Blob{}

	query := blobs.CreateBlob()
	query.StemWords(search_query)

	go func() {
		defer func() {
			bState.isLoading = false
		}()

		res := &blobs.BlobList{}
		var err error

		res, err = rep.UserMakeQuery(query)

		if err != nil || len(res.Blobs) == 0 {
			res, err = crw.CrawlIntoIndexer(search_query)
		}

		if res != nil && len(res.Blobs) > 0 {
			bState.items = res.Calculate_tf_idf(query)
			return
		}
		bState.items = []*blobs.Blob{}
	}()

	ti := textinput.New()
	ti.Placeholder = "Search again!"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 40
	ti.SetWidth(40)

	s := spinner.New()
	s.Spinner = spinner.Dot

	vi := viewport.New()
	return &results_screen{
		queryMade: search_query,

		searchInput: ti,
		viewport:    vi,

		spinner: s,
	}, s.Tick
}

var (
	CURRENT_SELECTOR        int     = 0
	HEADER_FOCUSEABLE_ITEMS int     = 1
	SCROLL_PERCENTAGE       float64 = 0.90
)

func (m *results_screen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if ok := m.searchInput.Focused(); ok {
				if len(m.searchInput.Value()) > MAX_CHAR_REQUIRED {
					return changeCurrentScreen(CreateResultsScreen(m.searchInput.Value()))
				}
			}

			if CURRENT_SELECTOR != 0 && len(bState.items) >= CURRENT_SELECTOR-1 {
				item := bState.items[CURRENT_SELECTOR-1]
				if item != nil {
					return changeCurrentScreen(CreateBlobInfoScreen(item, m.queryMade))
				}

			}

		case "down":
			if CURRENT_SELECTOR+1 < len(bState.items)+HEADER_FOCUSEABLE_ITEMS {
				CURRENT_SELECTOR += 1

				if m.GetScrollNeeded() > m.viewport.YOffset()+m.viewport.Height() {
					m.viewport.ScrollDown(m.viewport.Height())
				}

			}

		case "up":
			if CURRENT_SELECTOR-1 >= 0 {
				CURRENT_SELECTOR -= 1

				if m.GetScrollNeeded() < m.viewport.YOffset() {
					m.viewport.ScrollUp(m.viewport.Height())
				}

			}

		case "q":
			if bState.isLoading {
				return changeCurrentScreen(CreateMainScreen())
			}
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView(msg.Width))
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
	// m.viewport, cmd = m.viewport.Update(msg)
	// cmds = append(cmds, cmd)

	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	if bState.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

var (
	// general

	listMargin int = MARGIN_SIDES * 2
	list           = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			Padding(0, 1).
			Margin(1, listMargin).
			MaxWidth(150).
			MaxHeight(descriptionMaxHeight + infoUrlMaxHeight + headerMaxHeight + 2 + (listMargin * 2)) // 2 for padding

	// header
	headerMaxHeight int = 2
	headerTitle         = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			Bold(true).
			MaxHeight(headerMaxHeight)

	headerDate = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			Foreground(lipgloss.BrightBlack).
			MaxHeight(headerMaxHeight)

	header = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("#444")).
		MaxHeight(headerMaxHeight)

	// info card
	infoUrlMaxHeight int = 1
	infoUrl              = lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Left).
				Foreground(lipgloss.Blue).
				Underline(true).
				MaxHeight(infoUrlMaxHeight)

	infoScore = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right).
			MaxHeight(infoUrlMaxHeight)

	// description
	descriptionMaxHeight int = 2
	bottomDescription        = lipgloss.NewStyle().
				PaddingTop(1).
				Align(lipgloss.Left).
				MaxHeight(descriptionMaxHeight)

	// extras
	noItems = lipgloss.NewStyle().Italic(true).Bold(true).MarginTop(4).Align(lipgloss.Center)

	// center sections
	centerContainer = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Top)

	viewportItemsCenter = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Top)
)

func (m *results_screen) View(w, h int) tea.View {
	// render blobs info
	itemsListed := []string{}

	if !bState.isLoading && len(bState.items) <= 0 {
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

	var v tea.View
	if !m.ready {
		v.SetContent("\n  Initializing...")
	} else {

		m.viewport.SetContent(
			viewportItemsCenter.Width(w).Height(h).Render(itemListedRendered),
		)

		cc := lipgloss.NewStyle().Width(w)
		var container string = ""

		if bState.isLoading {
			spin := fmt.Sprintf("\n\n%s Loading forever...press q to quit\n\n", m.spinner.View())
			container = cc.Render(spin)
		} else {
			container = cc.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					m.headerView(w), m.viewport.View(),
				),
			)
		}

		v.SetContent(container)
	}

	return v
}

var titleSearch = lipgloss.NewStyle().Padding(0, 2).Margin(0, 20).Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.BrightBlack).AlignHorizontal(lipgloss.Left)

func (m *results_screen) headerView(w int) string {
	search := titleSearch.Width(w - 40) // 20 margin * 2 = 40
	if CURRENT_SELECTOR < HEADER_FOCUSEABLE_ITEMS {
		search = search.Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.Yellow)
		m.searchInput.Focus()
	} else {
		m.searchInput.Blur()
	}

	if !bState.isLoading {
		return search.Render(m.searchInput.View())
	}

	return ""
}

func (m *results_screen) bodyView(w, _ int) []string {
	itemsListed := []string{}
	for index, i := range bState.items {
		scoreParsed := int(i.Score * 100)

		list = list.
			BorderForeground(lipgloss.Red).
			Width(w - listMargin*2)

		if CURRENT_SELECTOR >= HEADER_FOCUSEABLE_ITEMS-1 && index == (CURRENT_SELECTOR-HEADER_FOCUSEABLE_ITEMS) {
			list = list.Border(lipgloss.DoubleBorder(), true).BorderForeground(lipgloss.BrightYellow)
		} else {
			list = list.Border(lipgloss.NormalBorder(), true)
		}

		// NOTE: header card (title + date)
		formattedDate := i.Datetime.String()
		titleUpper := strings.ToUpper(i.Title)
		headerTitleStr := headerTitle.
			Width(list.GetWidth() - len(formattedDate) - (listMargin * 2)).
			Render(titleUpper)

		headerDateStr := headerDate.Render(formattedDate)

		var headerStr lipgloss.Style = header
		if strings.Contains(titleUpper, "DISAMBIGUATION") == true {
			header = headerStr.Background(lipgloss.Red)
		} else {
			headerStr = headerStr.UnsetBackground()
		}

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
		bottomDescriptionStr := bottomDescription.Width(list.GetWidth()).Render(i.Description)

		// united all
		itemsListed = append(itemsListed,
			list.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					headerStr.Render(headerTitleStr, headerDateStr),
					informationCard,
					bottomDescriptionStr,
				),
			),
		)
	}
	return itemsListed
}

func (_ results_screen) GetScrollNeeded() int {
	return CURRENT_SELECTOR * ((list.GetHeight() * 2) + (listMargin * 2) + 2) // 2 for the padding
}
