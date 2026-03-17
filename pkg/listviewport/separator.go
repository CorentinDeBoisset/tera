package listviewport

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/google/uuid"
)

type SeparatorVariant int

const (
	SEPARATOR_BLANK SeparatorVariant = iota
	SEPARATOR_LINE
	SEPARATOR_DECORATED
)

type SeparatorModel struct {
	id           string
	variant      SeparatorVariant
	currentStyle lipgloss.Style
	width        int
}

func NewSeparator(width int, color color.Color, variant SeparatorVariant) *SeparatorModel {
	return &SeparatorModel{
		id:           uuid.NewString(),
		variant:      variant,
		currentStyle: lipgloss.NewStyle().Foreground(color).AlignHorizontal(lipgloss.Center).Width(width),
		width:        width,
	}
}

func (m *SeparatorModel) Resize(width int) {
	m.width = width
	m.currentStyle = m.currentStyle.Width(width)
}

func (m *SeparatorModel) Focusable() bool {
	return false
}

func (m *SeparatorModel) Height() int {
	return 1
}

func (m *SeparatorModel) View() string {
	switch m.variant {
	case SEPARATOR_DECORATED:
		var contentWidth int
		if m.width < 10 {
			return m.currentStyle.Render(strings.Repeat("─", contentWidth))
		} else if m.width < 30 {
			contentWidth = m.width * 60 / 100
		} else {
			contentWidth = m.width * 35 / 100
		}

		contentHalfWidth := (contentWidth - 3) / 2
		content := strings.Repeat("─", contentHalfWidth) + " ⟡ " + strings.Repeat("─", contentHalfWidth)
		return m.currentStyle.Render(content)

	case SEPARATOR_LINE:
		return m.currentStyle.Render(strings.Repeat("─", m.width*60/100))
	}

	return m.currentStyle.Render("")
}

func (s *SeparatorModel) Id() string {
	return s.id
}
