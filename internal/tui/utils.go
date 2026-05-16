package tui

import (
	"fmt"
	"image/color"
	"strings"

	"search_engine/internal/blobs"

	"charm.land/glamour/v2"
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

type HANDLE_DESCRIPTION struct {
	originalStr string // the entire blob buffer
	pages       [][]string
	index       int // points to page. renders x/2 chars after and x/2 chars before itself
}

var max_lines_per_page = 100

func (h *HANDLE_DESCRIPTION) returnPage() string {
	if h.index > h.getPagesLength() {
		return h.returnError("Out of bounds")
	}

	p := h.pages[h.index]
	str := strings.Join(p, "\n")

	if len(str) == 0 {
		return h.returnError("No text was supplied to the blob body")
	}

	l, err := glamour.Render(str, "dark")
	if err != nil {
		return str
	}

	return l
}

func (h *HANDLE_DESCRIPTION) returnError(reason string) string {
	t := fmt.Sprintf(`
	# Error!
	*Couldn't load the information! Sowwy!*
	Reason: %s
	`, reason)

	l, err := glamour.Render(t, "dark")
	if err != nil {
		return t
	}

	return l
}

// TODO:
// with help of the viewport...
// - make an initial render of X chars
// - when the viewport reachs near the bottom, render X amount more
// - override half of the previous string, ending with a "top" of max x/2 and a bottom to discover of x... x/2 + x in total
// - no lags! :)
func (h *HANDLE_DESCRIPTION) getDescriptionText(b *blobs.Blob) string {
	if len(h.pages) != 0 {
		return h.returnPage()
	}

	h.originalStr = b.GetBodyContent()

	// improve this...
	parsed := strings.ToValidUTF8(h.originalStr, "")
	lines := strings.Split(parsed, "\n")

	if len(lines) == 0 {
		return h.returnError("Invalid text provided")
	}

	p := []string{}

	for i, s := range lines {
		if i > 0 && i%max_lines_per_page == 0 {
			h.pages = append(h.pages, p)
			p = []string{}
			continue
		}

		p = append(p, strings.Join(strings.Fields(s), " "))
		if i == len(lines)-1 {
			h.pages = append(h.pages, p)
			break
		}
	}

	h.index = 0

	return h.returnPage()
}

// toBottom means thats is going to render the chars AFTER the index.
// yes, this needs to be made with a type string (enum). i dont car.
func (h *HANDLE_DESCRIPTION) resizeDescriptionText(nextPage bool) string {
	if nextPage {
		if h.index < h.getPagesLength() {
			h.index += 1
		}

		return h.returnPage()
	}

	if h.index > 0 {
		h.index -= 1
	}

	return h.returnPage()
}

func (h *HANDLE_DESCRIPTION) isAtBottom() bool {
	return h.index >= len(h.originalStr)
}

func (h *HANDLE_DESCRIPTION) getPagesLength() int {
	return len(h.pages) - 1
}
