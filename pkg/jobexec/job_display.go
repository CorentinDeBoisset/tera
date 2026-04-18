package jobexec

import (
	"fmt"
	"time"
	"unicode/utf8"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/corentindeboisset/tera/pkg/cfg"
	"github.com/corentindeboisset/tera/pkg/cmdrunr"
	"github.com/corentindeboisset/tera/pkg/iface"
	"github.com/corentindeboisset/tera/pkg/listviewport"
	"github.com/corentindeboisset/tera/pkg/outputviewer"
)

type RefreshStatusMsg time.Time

type keymap = struct {
	up, down, tab, quit key.Binding
}

type ifaceModel struct {
	width  int
	height int
	theme  iface.Theme

	keymap keymap
	help   help.Model

	jobConfig *cfg.JobConfig
	statuses  []StepStatus

	focusableTasks map[string]*JobItemView
	outputs        map[string]*cmdrunr.SafeBuffer

	stepPanelWidth  int
	focusOutput     bool
	hideOutputPanel bool
	focusedTask     string
	spinner         *iface.SharedSpinner
	stepPanel       listviewport.Model
	outputPanel     outputviewer.Model
}

func tickReadOutputsMsg() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return RefreshStatusMsg(t)
	})
}

func newModel(config *cfg.JobConfig, statuses []StepStatus, theme iface.Theme) ifaceModel {
	m := ifaceModel{
		help:  help.New(),
		theme: theme,
		keymap: keymap{
			up: key.NewBinding(
				key.WithKeys("k", "up"),
				key.WithHelp("↑/k", "Move up"),
			),
			down: key.NewBinding(
				key.WithKeys("j", "down"),
				key.WithHelp("↓/j", "Move down"),
			),
			tab: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("⇥/tab  ", "Switch focus"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("Ctrl+C ", "Exit"),
			),
		},
		focusOutput:     false,
		hideOutputPanel: false,
		jobConfig:       config,
		stepPanelWidth:  15,
		statuses:        statuses,
		spinner:         iface.NewSharedSpinner(),
		stepPanel:       listviewport.New(30, 10, lipgloss.NewStyle()),
		outputPanel:     outputviewer.New(30, 10, theme, nil),
	}

	m.updateKeyBindings()
	m.calculateMinPanelSize()
	m.outputPanel.SetFocus(false)

	m.initializeStepPanel()

	return m
}

func (m *ifaceModel) initializeStepPanel() {
	jobItemViews := make([]listviewport.ListItem, 0)
	m.focusableTasks = make(map[string]*JobItemView)
	m.outputs = make(map[string]*cmdrunr.SafeBuffer)

	for stepIdx := range m.statuses {
		id := fmt.Sprintf("step#%d", stepIdx)
		jobItemViews = append(jobItemViews, NewJobItemView(id, false, 0, m.jobConfig.Steps[stepIdx].Name, &m.statuses[stepIdx].Stater, m.spinner, m.theme, m.width))

		if m.statuses[stepIdx].BeforeHooks != nil {
			id := fmt.Sprintf("step#%d__bh", stepIdx)
			m.focusableTasks[id] = NewJobItemView(id, true, 2, "Pre-run hooks", &m.statuses[stepIdx].BeforeHooks.Stater, m.spinner, m.theme, m.width)
			m.outputs[id] = &m.statuses[stepIdx].BeforeHooks.Output
			jobItemViews = append(jobItemViews, m.focusableTasks[id])
		}

		for taskIdx := range m.statuses[stepIdx].Tasks {
			id := fmt.Sprintf("step#%d__task#%d", stepIdx, taskIdx)
			m.focusableTasks[id] = NewJobItemView(id, true, 2, m.jobConfig.Steps[stepIdx].Tasks[taskIdx].Name, &m.statuses[stepIdx].Tasks[taskIdx].Stater, m.spinner, m.theme, m.width)
			m.outputs[id] = &m.statuses[stepIdx].Tasks[taskIdx].Output
			jobItemViews = append(jobItemViews, m.focusableTasks[id])
		}

		if m.statuses[stepIdx].AfterHooks != nil {
			id := fmt.Sprintf("step#%d__ah", stepIdx)
			m.focusableTasks[id] = NewJobItemView(id, true, 2, "Post-run hooks", &m.statuses[stepIdx].AfterHooks.Stater, m.spinner, m.theme, m.width)
			m.outputs[id] = &m.statuses[stepIdx].AfterHooks.Output
			jobItemViews = append(jobItemViews, m.focusableTasks[id])
		}

		if stepIdx < len(m.statuses)-1 {
			jobItemViews = append(jobItemViews, listviewport.NewSeparator(m.width, lipgloss.Black, listviewport.SEPARATOR_BLANK))
		}
	}

	m.stepPanel.SetItems(jobItemViews)

	// Initialize the focus
	focusedId := m.stepPanel.GoToTop()
	m.updateFocus(focusedId)
}

func (m *ifaceModel) calculateMinPanelSize() {
	maxLen := 10
	for _, step := range m.jobConfig.Steps {
		if maxLen < utf8.RuneCountInString(step.Name) {
			maxLen = utf8.RuneCountInString(step.Name)
		}
		for _, task := range step.Tasks {
			if maxLen < utf8.RuneCountInString(task.Name)+2 {
				maxLen = utf8.RuneCountInString(task.Name) + 2
			}
		}
		if len(step.RunBefore) > 0 && maxLen < 19 {
			maxLen = 19
		}
		if len(step.RunAfter) > 0 && maxLen < 18 {
			maxLen = 18
		}
	}

	m.stepPanelWidth = maxLen + 10
}

func (m *ifaceModel) updateSizes() {
	// The height is fixed
	panelsHeight := m.height - 5

	stepPanelWidth := m.stepPanelWidth
	if stepPanelWidth > (m.width/2 - 6) {
		stepPanelWidth = m.width/2 - 6
	}

	outputWidth := m.width - stepPanelWidth
	if outputWidth < 10 {
		m.hideOutputPanel = true
		m.stepPanel.Resize(m.width, panelsHeight)
		return
	}

	m.outputPanel.Resize(outputWidth, panelsHeight)
	m.stepPanel.Resize(stepPanelWidth, panelsHeight)
}

func (m *ifaceModel) updateKeyBindings() {
	m.keymap.tab.SetEnabled(!m.hideOutputPanel)
}

func (m ifaceModel) Init() tea.Cmd {
	return tea.Batch(
		tickReadOutputsMsg(),
		m.spinner.Tick,
	)
}

func (m *ifaceModel) updateFocus(focusedId string) {
	if _, ok := m.focusableTasks[m.focusedTask]; ok {
		m.focusableTasks[m.focusedTask].SetFocus(0)
	}

	m.focusedTask = focusedId

	if _, ok := m.focusableTasks[focusedId]; ok {
		m.focusableTasks[focusedId].SetFocus(2)
		m.outputPanel.SetBuffer(m.outputs[focusedId])
	}
}

func (m ifaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Global (independant of the panel with focus)
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab":
			if !m.hideOutputPanel {
				m.focusOutput = !m.focusOutput
				m.outputPanel.SetFocus(m.focusOutput)
				m.updateKeyBindings()
				if m.focusOutput {
					m.focusableTasks[m.focusedTask].SetFocus(1)
				} else {
					m.focusableTasks[m.focusedTask].SetFocus(2)
				}
			}
		}

		if m.focusOutput {
			return m, m.outputPanel.Update(msg)
		}

		switch msg.String() {
		case "up", "k":
			// TODO: create a `CircularScrollUp(n int) string` method
			focusedId := m.stepPanel.ScrollUp(1)
			m.updateFocus(focusedId)

		case "down", "j":
			// TODO: create a `CircularScrollDown(n int) string` method
			focusedId := m.stepPanel.ScrollDown(1)
			m.updateFocus(focusedId)

		// Other movement keys, not displayed in the help
		case "pgup":
			focusedId := m.stepPanel.PageUp()
			m.updateFocus(focusedId)

		case "pgdown":
			focusedId := m.stepPanel.PageDown()
			m.updateFocus(focusedId)

		case "home":
			focusedId := m.stepPanel.GoToTop()
			m.updateFocus(focusedId)

		case "end":
			focusedId := m.stepPanel.GoToBottom()
			m.updateFocus(focusedId)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

		// Clear the screen to avoid artifacts
		return m, tea.ClearScreen

	case RefreshStatusMsg:
		if !m.hideOutputPanel {
			m.outputPanel.RefreshContent()
		}
		return m, tickReadOutputsMsg()

	case tea.PasteMsg:
		if m.focusOutput {
			return m, m.outputPanel.Update(msg)
		}

	case spinner.TickMsg:
		return m, m.spinner.Update(msg)

	case tea.MouseWheelMsg:
		if m.focusOutput {
			return m, m.outputPanel.Update(msg)
		}

		switch msg.Button {
		case tea.MouseWheelUp:
			focusedId := m.stepPanel.ScrollUp(1)
			m.updateFocus(focusedId)

		case tea.MouseWheelDown:
			focusedId := m.stepPanel.ScrollDown(1)
			m.updateFocus(focusedId)
		}
	}

	return m, nil
}

func (m ifaceModel) View() tea.View {
	help := m.help.FullHelpView([][]key.Binding{
		{m.keymap.up, m.keymap.down},
		{m.keymap.tab, m.keymap.quit},
	})

	var views []string
	views = append(views, m.stepPanel.View())
	if !m.hideOutputPanel {
		views = append(views, m.outputPanel.View())
	}

	view := tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help)
	view.AltScreen = true
	view.KeyboardEnhancements.ReportEventTypes = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}
