// Package overlay provides simple overlay windows
package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Anchor int

const (
	Begin Anchor = iota
	Center
	End
)

// Model contains any tea.Model inside
type Model struct {
	Content  tea.Model
	IsActive bool   // Flag for upper composition (Default: `false`)
	Width    int    // Set to positive if you want specific width (Default: `0`)
	Height   int    // Set to positive if you want specific height (Default: `0`)
	XAnchor  Anchor // Start position (Default: `Begin` - very top)
	YAnchor  Anchor // Start position (Default: `Begin` - very left)
	XOffset  int    // Counts from the `XAnchor` (Default: `0`)
	YOffset  int    // Counts from the `YAnchor` (Default: `0`)
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
	if m.Content == nil {
		return ""
	}
	layout := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#ffffff"))
	if m.Width > 0 {
		layout = layout.Width(m.Width)
	}
	if m.Height > 0 {
		layout = layout.Height(m.Height)
	}
	return layout.Render(m.Content.View())
}

func New(content tea.Model) *Model {
	return &Model{
		Content: content,
	}
}

func (m *Model) Place(bg string) string {
	return Compose(m.View(), bg, m.XAnchor, m.YAnchor, m.XOffset, m.YOffset)
}
