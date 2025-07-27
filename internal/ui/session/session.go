// Package session contains Model, which represents main window of TUI
package session

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/popup"
	"github.com/alphameo/nm-tui/internal/ui/styles"
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
	state        sessionState
	wifi         wifi.Model
	timer        timer.Model
	floatWin     popup.Model
	notification popup.Model
	width        int
	height       int
}

func New() Model {
	w := wifi.New()
	t := timer.New(time.Hour)
	f := popup.New(NewTextModel())
	n := popup.New(NewTextModel())
	m := Model{
		wifi:         w,
		timer:        t,
		floatWin:     f,
		notification: n,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		m.wifi.Init(),
		m.floatWin.Init(),
		m.notification.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	var upd tea.Model
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(msg)
		m.notification = upd.(popup.Model)
		cmds = append(cmds, cmd)
	} else if m.floatWin.IsActive {
		upd, cmd = m.floatWin.Update(msg)
		m.floatWin = upd.(popup.Model)
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
				m.floatWin.IsActive = true
			case "n":
				m.notify("xdd")
			}
		}
	}
	size, ok := msg.(tea.WindowSizeMsg)
	if ok {
		m.width = size.Width - 2
		m.height = size.Height - 2
	}
	upd, cmd = m.wifi.Update(msg)
	cmds = append(cmds, cmd)
	m.wifi = upd.(wifi.Model)
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	mainView := m.wifi.View() + "\n" + m.timer.View() + fmt.Sprintf("\n state: %v", m.state)
	if m.floatWin.IsActive {
		popupLayout := lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			m.floatWin.View(),
		)
		mainView = popupLayout
	}
	if m.notification.IsActive {
		notifLayout := lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			m.notification.View(),
		)
		mainView = notifLayout
	}
	return styles.BorderStyle.Width(m.width).Height(m.height).Render(mainView)
}

func (m *Model) showPopup(content tea.Model) {
	m.floatWin.IsActive = true
	m.floatWin.Content = content
}

func (m *Model) notify(text string) {
	t, ok := m.notification.Content.(TextModel)
	if !ok {
		logger.ErrorLog.Println("Invalid Type")
	}
	t.Text = text
	m.notification.Content = t
	m.notification.IsActive = true
}

// TextModel contains only string inside
type TextModel struct{ Text string }

func (m TextModel) Init() tea.Cmd {
	return nil
}

func (m TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TextModel) View() string {
	return m.Text
}

func NewTextModel() TextModel {
	return TextModel{"Placeholder"}
}
