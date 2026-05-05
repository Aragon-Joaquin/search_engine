package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

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
