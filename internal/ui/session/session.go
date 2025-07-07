package session

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/ui/wifi"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	wifiView sessionState = iota
	timerView
)

type Model struct {
	state sessionState
	wifi  wifi.Model
	timer timer.Model
}

func New() Model {
	wifi := wifi.New()
	timer := timer.New(time.Hour)
	m := Model{wifi: wifi, timer: timer}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), m.wifi.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.state == wifiView {
				m.state = timerView
			} else {
				m.state = wifiView
			}
		}
		switch m.state {
		case wifiView:
			var updated tea.Model
			updated, cmd = m.wifi.Update(msg)
			cmds = append(cmds, cmd)
			m.wifi = updated.(wifi.Model)
		}
	default:
		var updated tea.Model
		updated, cmd = m.wifi.Update(msg)
		cmds = append(cmds, cmd)
		m.wifi = updated.(wifi.Model)
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.wifi.View() + "\n" + m.timer.View() + fmt.Sprint(m.state)
}
