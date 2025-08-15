// Package label provides simple model with text, which should be shown to user
package label

import tea "github.com/charmbracelet/bubbletea"

type Model string

func New(label string) Model {
	return Model(label)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return string(m)
}
