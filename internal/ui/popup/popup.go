package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Content string
}

type CloseMsg bool

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, m.quit()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	overlay := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(30).
		Height(7).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#ffffff"))
	return overlay.Render(m.Content)
}

func New() Model {
	return Model{}
}

func (m *Model) quit() tea.Cmd {
	return func() tea.Msg {
		return CloseMsg(true)
	}
}
