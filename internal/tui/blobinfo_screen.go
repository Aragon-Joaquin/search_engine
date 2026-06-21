package tui

import (
	"fmt"
	"strings"

	"search_engine/internal/blobs"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var description_state HANDLE_DESCRIPTION

type blobinfo_screen struct {
	blob      *blobs.Blob
	prevQuery string

	viewport viewport.Model
}

func CreateBlobInfoScreen(b *blobs.Blob, pq string) CurrentScreen {
	description_state = HANDLE_DESCRIPTION{
		originalStr: "",
		pages:       [][]string{},
		index:       0,
	}

	return &blobinfo_screen{
		blob:      b,
		prevQuery: pq,
		viewport:  viewport.New(),
	}
}

func (m *blobinfo_screen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			return changeCurrentScreen(CreateResultsScreen(m.prevQuery))

		case "left":
			description_state.resizeDescriptionText(false)

		case "right":
			description_state.resizeDescriptionText(true)

		}

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height)
		m.viewport.YPosition = buttonBack.GetWidth()
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

var (
	title = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Red).
		Bold(true).
		Padding(1).
		MarginBottom(1)

	description = lipgloss.NewStyle().
			PaddingBottom(8).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Top)

	pageHelper = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(lipgloss.Red).
			Bold(true).
			Height(2)
)

func (m *blobinfo_screen) View(w, h int) tea.View {
	marginWidth := int(float64(w) / 1.3)
	// body header
	titleStr := title.MaxWidth(MAX_WIDTH_LIST).Width(marginWidth).Render(strings.ToUpper(m.blob.Title))
	titleStr = lipgloss.Place(w, 0, lipgloss.Center, lipgloss.Center, titleStr)

	// body body
	d := description_state.getDescriptionText(m.blob)

	descriptioSub := m.createSubtitle("Description", marginWidth)
	descriptionStr := description.Render(d)

	var navPage string
	if description_state.index < description_state.getPagesLength() {
		navPage = pageHelper.Width(w).Render("<--- View More --->")
	} else {
		navPage = pageHelper.Width(w).Render("END OF CONTENT")
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStr,
			navPage,
			lipgloss.
				Place(w, h, lipgloss.Center, lipgloss.Top,
					fmt.Sprintf("%s\n%s", descriptioSub, descriptionStr),
				),
		),
	)

	var v tea.View
	v.SetContent(fmt.Sprintf("%s\n%s",
		m.headerView(h, w),
		m.viewport.View(),
	),
	)

	return v
}

var (
	buttonText = "<- Back"
	buttonBack = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(lipgloss.Yellow).
			Padding(0, 2).
			Width(len(buttonText) + 3*2)

	urlSpan = lipgloss.NewStyle().
		Foreground(lipgloss.Blue).
		Align(lipgloss.Center).
		Margin(1, 0).
		Underline(true)

	pagecountSpan = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Margin(1, 0).
			Underline(true).
			Foreground(lipgloss.BrightBlack)

	borderHeader = lipgloss.NewStyle().
			Border(lipgloss.ASCIIBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#222"))
)

// NOTE: private
func (m *blobinfo_screen) headerView(_, w int) string {
	text := fmt.Sprintf("Page %d out of %d", description_state.index, description_state.getPagesLength())

	pageCount := pagecountSpan.
		Width(len(text) + 4)
	goBack := buttonBack.Render(buttonText)

	url := urlSpan.
		Hyperlink(m.blob.URL).
		Height(buttonBack.GetHeight()).
		Width(w - (buttonBack.GetWidth() + pageCount.GetWidth())).
		Render(m.blob.URL)
	return borderHeader.Width(w).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		goBack,
		url,
		pageCount.Render(text),
	))
}

var subtitle = lipgloss.NewStyle().Foreground(lipgloss.Magenta).Bold(true).MarginBottom(1)

func (m *blobinfo_screen) createSubtitle(titlename string, w int) string {
	return subtitle.Width(w).MaxWidth(MAX_WIDTH_LIST).AlignHorizontal(lipgloss.Left).Render(fmt.Sprintf("## %s", titlename))
}
