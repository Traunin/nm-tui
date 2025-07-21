package popup

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	content string
	Active  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Active = false
			return m, nil
		}
	}
	return m, nil
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

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
