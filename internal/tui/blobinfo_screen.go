package tui

import (
	"fmt"

	"search_engine/internal/blobs"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type blobinfo_screen struct {
	blob      *blobs.Blob
	prevQuery string

	viewport viewport.Model
}

func CreateBlobInfoScreen(b *blobs.Blob, pq string) CurrentScreen {
	return &blobinfo_screen{
		blob:      b,
		prevQuery: pq,
		viewport:  viewport.New(),
	}
}

func (m *blobinfo_screen) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			return changeCurrentScreen(CreateResultsScreen(m.prevQuery))
		}

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - buttonBack.GetWidth())

		m.viewport.YPosition = buttonBack.GetWidth()
	}
	return cmd
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
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Top)
)

func (m *blobinfo_screen) View(w, h int) tea.View {
	marginWidth := int(float64(w) / 1.3)
	// body header
	titleStr := title.Width(marginWidth).Render(m.blob.Title)
	titleStr = lipgloss.Place(w, 0, lipgloss.Center, lipgloss.Center, titleStr)

	// body body
	descriptionStr := description.Width(marginWidth + 8).Render(m.blob.Description)
	descriptioSub := m.createSubtitle("Description", marginWidth+8)

	descriptionStr = lipgloss.Place(w, h, lipgloss.Center, lipgloss.Top, fmt.Sprintf("%s\n%s", descriptioSub, descriptionStr))

	m.viewport.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStr,
			descriptionStr,
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

	borderHeader = lipgloss.NewStyle().
			Border(lipgloss.ASCIIBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#222"))
)

func (m *blobinfo_screen) headerView(_, w int) string {
	goBack := buttonBack.Render(buttonText)
	url := urlSpan.
		Hyperlink(m.blob.URL).
		Height(buttonBack.GetHeight()).
		Width(w - (len(m.blob.URL) / 2) - (buttonBack.GetWidth() / 2) - 1).
		Render(m.blob.URL)

	return borderHeader.Width(w).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		goBack,
		url,
	))
}

var subtitle = lipgloss.NewStyle().Foreground(lipgloss.Magenta).Bold(true).MarginBottom(1)

func (m *blobinfo_screen) createSubtitle(titlename string, w int) string {
	return subtitle.Width(w).Render(fmt.Sprintf("## %s", titlename))
}
