package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	content string
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
		case "esc":
			return m, m.quit()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	overlay := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(30).
		Height(7).
		Align(lipgloss.Center).
		Background(lipgloss.Color("#111111")).
		Foreground(lipgloss.Color("#ffffff")).
		Render("Floating window!\nPress 'o' to close.")
	return overlay
}

func New() Model {
	return Model{}
}

func (m *Model) quit() tea.Cmd {
	return func() tea.Msg {
		return CloseMsg(true)
	}
}
