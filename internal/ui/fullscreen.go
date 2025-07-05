package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model int

type tickMsg time.Time

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick()
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("\n\n OHIO in %d", m)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
