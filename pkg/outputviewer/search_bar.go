package outputviewer

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const search_prefix = "Search: "

func sanitize(input string) []rune {
	sanitized := make([]rune, 0, 2*len(input)) // worst case scenario: all characters must be escaped

	for idx, r := range input {
		switch r {
		case '\t':
			sanitized = append(sanitized, '\\', 't')
		case '\r':
			// If followed by \n, we just skip the \r
			if idx < len(input)-1 && input[idx+1] == '\n' {
				continue
			}

			sanitized = append(sanitized, '\\', 'n')
		case '\n':
			sanitized = append(sanitized, '\\', 'n')
		default:
			if unicode.IsControl(r) {
				continue
			}
			sanitized = append(sanitized, r)
		}
	}

	return sanitized
}

type SearchBarModel struct {
	width        int
	contentWidth int

	currentSearch []rune
	cursorIdx     int
	scrollIdx     int
	showCursor    bool
	virtualCursor cursor.Model
}

func NewSearchBar(width int) *SearchBarModel {
	model := &SearchBarModel{
		width:         10,
		contentWidth:  10,
		currentSearch: nil,
		cursorIdx:     0,
		scrollIdx:     0,
		showCursor:    false,
		virtualCursor: cursor.Model{},
	}

	// Recompute the width and contentWidth
	model.Resize(width)

	return model
}

func (m *SearchBarModel) Resize(width int) {
	m.width = width

	// The -1 accounts for an empty space at the end for the cursor
	m.contentWidth = width - len(search_prefix) - 1

	// No-op repositioning of the cursor, to recompute the scroll
	m.moveCursor(0)
}

func (m *SearchBarModel) scroll(pos int) {
	m.scrollIdx = max(min(pos, len(m.currentSearch)-m.contentWidth), 0)
}

func (m *SearchBarModel) insertRunes(input []rune) {
	m.currentSearch = slices.Concat(m.currentSearch[:m.cursorIdx], input, m.currentSearch[m.cursorIdx:])
	m.moveCursor(len(input))
}

// Move the cursor to the absolute index `idx`
func (m *SearchBarModel) moveCursorAbs(idx int) {
	newIndex := min(max(idx, 0), len(m.currentSearch))
	m.cursorIdx = newIndex
	if m.cursorIdx == len(m.currentSearch) {
		m.virtualCursor.SetChar(" ")
		m.scroll(m.cursorIdx - m.contentWidth)
	} else {
		m.virtualCursor.SetChar(string(m.currentSearch[m.cursorIdx]))
		if newIndex < m.scrollIdx {
			// Scroll to the left
			m.scroll(newIndex)
		} else if (newIndex - m.contentWidth) > m.scrollIdx {
			// Scroll to the right
			m.scroll(newIndex - m.contentWidth)
		}
	}
}

// Move the cursor of `n` relative to the current position
func (m *SearchBarModel) moveCursor(n int) {
	m.moveCursorAbs(m.cursorIdx + n)
}

func (m *SearchBarModel) Submit() *regexp.Regexp {
	if len(m.currentSearch) == 0 {
		return nil
	}

	reg, err := regexp.Compile(string(m.currentSearch))
	if err != nil {
		log.Printf("The requested regexp is invalid: %s", err)
		return nil
	}

	return reg
}

func (m *SearchBarModel) HandleKeyMsg(msg tea.KeyPressMsg) {
	switch msg.String() {
	case "left":
		m.moveCursor(-1)
	case "right":
		m.moveCursor(1)
	case "home":
		m.moveCursorAbs(0)
	case "end":
		m.moveCursorAbs(len(m.currentSearch))
	case "backspace":
		if len(m.currentSearch) > 0 && m.cursorIdx > 0 {
			m.currentSearch = slices.Concat(m.currentSearch[:max(0, m.cursorIdx-1)], m.currentSearch[m.cursorIdx:])
			m.moveCursor(-1)
		}
	case "delete":
		if len(m.currentSearch) > 0 && m.cursorIdx < len(m.currentSearch) {
			m.currentSearch = slices.Concat(m.currentSearch[:m.cursorIdx], m.currentSearch[m.cursorIdx+1:])
			m.moveCursor(0)
		}
	default:
		if len(msg.Text) > 0 {
			// Important to sanitize, in case of a ctrl-v there are a lot of control characters (new lines, tabs...)
			m.InsertText(msg.Text)
		}
	}
}

func (m *SearchBarModel) InsertText(input string) {
	m.insertRunes(sanitize(input))
}

func (m *SearchBarModel) ToggleCursor(visible bool) {
	m.showCursor = visible
	m.moveCursorAbs(len(m.currentSearch))
}

func (m *SearchBarModel) SetCursorVisibility(focus bool) {
	m.virtualCursor.IsBlinked = !focus
}

func (m *SearchBarModel) Clear() {
	m.currentSearch = nil
	m.moveCursorAbs(0)
}

func (m *SearchBarModel) View() string {
	var content string
	startIdx := m.scrollIdx

	// The +1 accounts for an empty space at the end for the cursor
	endIdx := min(m.scrollIdx+m.contentWidth+1, len(m.currentSearch))

	if !m.showCursor {
		content = string(m.currentSearch[startIdx:endIdx])
	} else {
		parts := make([]string, 0, 3)
		if startIdx < m.cursorIdx {
			parts = append(parts, string(m.currentSearch[startIdx:m.cursorIdx]))
		}
		parts = append(parts, m.virtualCursor.View())
		if m.cursorIdx < endIdx {
			parts = append(parts, string(m.currentSearch[m.cursorIdx+1:endIdx]))
		}

		content = strings.Join(parts, "")
	}

	return lipgloss.NewStyle().
		Width(m.width).
		MaxWidth(m.width).
		Height(1).
		Render(fmt.Sprintf("%s%s", search_prefix, content))
}
