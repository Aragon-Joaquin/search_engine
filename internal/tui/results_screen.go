package tui

import (
	"math"
	"strconv"

	"search_engine/internal/blobs"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type results_screen struct {
	items []*blobs.Blob
}

func CreateResultsScreen(search_query string) CurrentScreen {
	return &results_screen{
		items: rep.UserMakeQuery(search_query),
	}
}

func (m *results_screen) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	return cmd
}

func (m *results_screen) View(w, h int) tea.View {
	itemsListed := []string{}

	for _, i := range m.items {
		formattedDate := i.Datetime.Format("2006/01/2")
		marginWidth := 4

		list := lipgloss.NewStyle().AlignVertical(lipgloss.Center).Margin(0, marginWidth).Width(w - (marginWidth * 2))
		halfWidth := int(math.Floor(float64(list.GetWidth() / 2)))

		var descriptionFallback string = "No description provided"
		if i.Description != "" {
			descriptionFallback = i.Description
		}

		scoreToStr := strconv.Itoa(int(i.Score))

		itemsListed = append(itemsListed,
			list.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.JoinHorizontal(
						lipgloss.Center,
						lipgloss.NewStyle().Width(halfWidth).AlignHorizontal(lipgloss.Left).Render(i.Title),
						lipgloss.NewStyle().Width(halfWidth).AlignHorizontal(lipgloss.Right).Render(formattedDate),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Center,
						lipgloss.NewStyle().Width(halfWidth).AlignHorizontal(lipgloss.Left).Render(i.URL),
						lipgloss.NewStyle().Width(halfWidth).AlignHorizontal(lipgloss.Right).Render("Score: "+scoreToStr),
					),
					descriptionFallback,
				),
			),
		)
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Center, itemsListed...))

	v.AltScreen = true
	return v
}
