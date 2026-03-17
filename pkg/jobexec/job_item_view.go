package jobexec

import (
	"fmt"
	"sync"

	"charm.land/lipgloss/v2"
	"github.com/corentindeboisset/tera/pkg/iface"
)

var successFlag = sync.OnceValue(func() string {
	return lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Bold(true).Render("✓")
})

var failureFlag = sync.OnceValue(func() string {
	return lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true).Render("✗")
})

type JobItemView struct {
	id        string
	focusable bool

	label   string
	stater  *Stater
	spinner *iface.SharedSpinner
	theme   iface.Theme

	focused bool

	currentStyle lipgloss.Style
	cachedHeight int
}

func NewJobItemView(id string, focusable bool, paddingLeft int, label string, stater *Stater, spinner *iface.SharedSpinner, theme iface.Theme, width int) *JobItemView {
	model := JobItemView{
		id:        id,
		focusable: focusable,

		label:   label,
		stater:  stater,
		spinner: spinner,
		theme:   theme,

		focused: false,

		currentStyle: lipgloss.NewStyle().PaddingLeft(paddingLeft),
		cachedHeight: 1,
	}
	model.Resize(width)

	return &model
}

func (m *JobItemView) Id() string {
	return m.id
}

func (m *JobItemView) Height() int {
	return m.cachedHeight
}

func (m *JobItemView) View() string {
	var content string

	labelStyle := lipgloss.NewStyle()
	if m.focused {
		labelStyle = labelStyle.Bold(true).Foreground(lipgloss.Blue).Underline(true)
	}

	switch m.stater.State() {
	case STATE_RUNNING:
		content = fmt.Sprintf("%s %s", m.spinner.View(), labelStyle.Render(m.label))
	case STATE_SUCCESSFUL:
		content = fmt.Sprintf("%s  %s", successFlag(), labelStyle.Render(m.label))
	case STATE_FAILED:
		content = fmt.Sprintf("%s  %s", failureFlag(), labelStyle.Render(m.label))
	case STATE_NOT_STARTED:
		fallthrough
	default:
		content = fmt.Sprintf("   %s", labelStyle.Render(m.label))
	}

	return m.currentStyle.Render(content)
}

func (m *JobItemView) Focusable() bool {
	return m.focusable
}

func (m *JobItemView) SetFocus(focus bool) {
	m.focused = m.focusable && focus
}

func (m *JobItemView) SetLabel(label string) {
	m.label = label
	m.cachedHeight = lipgloss.Height(m.View())
}

func (m *JobItemView) Resize(w int) {
	m.currentStyle = m.currentStyle.Width(w)
	m.cachedHeight = lipgloss.Height(m.View())
}
