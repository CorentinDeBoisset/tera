package iface

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SharedSpinner struct {
	spinner spinner.Model
}

func NewSharedSpinner() *SharedSpinner {
	return &SharedSpinner{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.BrightYellow)),
		),
	}
}

func (m *SharedSpinner) Update(msg tea.Msg) tea.Cmd {
	var newCmd tea.Cmd
	m.spinner, newCmd = m.spinner.Update(msg)

	return newCmd
}

func (m *SharedSpinner) View() string {
	return m.spinner.View()
}

func (m *SharedSpinner) Tick() tea.Msg {
	return m.spinner.Tick()
}
