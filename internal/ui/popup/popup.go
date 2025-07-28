// Package popup provides simple popup windows
package popup

import (
	"github.com/alphameo/nm-tui/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains any tea.Model inside
type Model struct {
	Content  tea.Model
	IsActive bool
	width    int
	height   int
}

func (m Model) Init() tea.Cmd {
	return m.Content.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.IsActive = false
			return m, nil
		}
	}
	m.Content, cmd = m.Content.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	overlay := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#ffffff"))
	return overlay.Render(m.Content.View())
}

func New(content tea.Model, width int, heigh int) Model {
	if content == nil {
		logger.ErrorLog.Panicln("content is nil")
	}
	return Model{content, false, width, heigh}
}
