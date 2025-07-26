package session

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/ui/popup"
	"github.com/alphameo/nm-tui/internal/ui/wifi"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	wifiView sessionState = iota
	timerView
)

type Model struct {
	state     sessionState
	wifi      wifi.Model
	timer     timer.Model
	popup     popup.Model
	popActive bool
}

func New() Model {
	wifi := wifi.New()
	timer := timer.New(time.Hour)
	popup := popup.New()
	m := Model{
		wifi:      wifi,
		timer:     timer,
		popup:     popup,
		popActive: false,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		m.wifi.Init(),
		m.popup.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.popActive {
		var upd tea.Model
		_, ok := msg.(popup.CloseMsg)
		if ok {
			m.popActive = !m.popActive
		}
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(popup.Model)
		cmds = append(cmds, cmd)
	} else {
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
			case "o":
				m.popActive = !m.popActive
			}
		}
	}
	var updated tea.Model
	updated, cmd = m.wifi.Update(msg)
	cmds = append(cmds, cmd)
	m.wifi = updated.(wifi.Model)
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	mainView := m.wifi.View() + "\n" + m.timer.View() + fmt.Sprint(m.state)

		return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center,
	if m.popActive {
			mainView+"\n"+m.popup.View())
	}
	return mainView
}
