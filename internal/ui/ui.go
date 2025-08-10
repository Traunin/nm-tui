// Package ui contains Model, which represents main window of TUI
package ui

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/ui/label"
	"github.com/alphameo/nm-tui/internal/ui/overlay"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	wifiView sessionState = iota
	timerView
)

type Model struct {
	state        sessionState
	wifiTable    WifiTableModel
	timer        timer.Model
	popup        overlay.Model
	notification overlay.Model
	width        int
	height       int
}

func New() Model {
	w := NewWifiTableModel(51, 20)
	t := timer.New(time.Hour)
	p := overlay.New(nil)
	p.Width = 100
	p.Height = 10
	p.XAnchor = overlay.Center
	p.YAnchor = overlay.Center
	n := overlay.New(nil)
	n.XAnchor = overlay.Center
	n.YAnchor = overlay.Center
	n.Width = 100
	n.Height = 10
	m := Model{
		wifiTable:    *w,
		timer:        t,
		popup:        *p,
		notification: *n,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		m.wifiTable.Init(),
		m.popup.Init(),
		m.notification.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case PopupContentMsg:
		return m, m.showPopup(msg)
	case PopupActivityMsg:
		m.popup.IsActive = bool(msg)
		return m, nil
	case NotificationTextMsg:
		return m, m.showNotification(string(msg))
	case NotificationActivityMsg:
		m.notification.IsActive = bool(msg)
		return m, nil
	case tea.KeyMsg:
		return m, m.processKeyMsg(msg)
	}
	return m, m.processCommonMsg(msg)
}

func (m Model) View() string {
	mainView := m.wifiTable.View() + "\n" + m.timer.View() + fmt.Sprintf("\n state: %v", m.state)
	mainView = styles.BorderStyle.Width(m.width - 2).Height(m.height - 2).Render(mainView)

	if m.popup.IsActive {
		mainView = m.popup.Place(mainView)
	}
	if m.notification.IsActive {
		mainView = m.notification.Place(mainView)
	}
	return mainView
}

func (m *Model) processKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(keyMsg)
		m.notification = upd.(overlay.Model)
		return cmd
	} else if m.popup.IsActive {
		upd, cmd = m.popup.Update(keyMsg)
		m.popup = upd.(overlay.Model)
		return cmd
	}
	switch keyMsg.String() {
	case "q", "ctrl+q", "esc", "ctrl+c":
		return tea.Quit
	case "tab":
		if m.state == wifiView {
			m.state = timerView
		} else {
			m.state = wifiView
		}
		return cmd
	}
	upd, cmd = m.wifiTable.Update(keyMsg)
	m.wifiTable = upd.(WifiTableModel)
	return cmd
}

func (m *Model) processCommonMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	m.timer, cmd = m.timer.Update(msg)
	if cmd != nil {
		return cmd
	}
	upd, cmd = m.wifiTable.Update(msg)
	m.wifiTable = upd.(WifiTableModel)
	if cmd != nil {
		return cmd
	}
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(msg)
		m.notification = upd.(overlay.Model)
		if cmd != nil {
			return cmd
		}
	}
	if m.popup.IsActive {
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(overlay.Model)
		if cmd != nil {
			return cmd
		}
	}
	return nil
}

func (m *Model) showPopup(content tea.Model) tea.Cmd {
	m.popup.IsActive = true
	m.popup.Content = content
	return m.popup.Content.Init()
}

func (m *Model) showNotification(text string) tea.Cmd {
	m.notification.IsActive = true
	m.notification.Content = label.New(text)
	return nil
}

// Public controls

type (
	PopupContentMsg  tea.Model
	PopupActivityMsg bool
)

func ShowPopup(content tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg(content)
	}
}

func ClosePopup() tea.Msg {
	return PopupActivityMsg(false)
}

type (
	NotificationTextMsg     string
	NotificationActivityMsg bool
)

func ShowNotification(notification string) tea.Cmd {
	return func() tea.Msg {
		return NotificationTextMsg(notification)
	}
}

func CloseNotification() tea.Msg {
	return NotificationActivityMsg(false)
}
