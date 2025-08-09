// Package ui contains Model, which represents main window of TUI
package ui

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/ui/label"
	"github.com/alphameo/nm-tui/internal/ui/overlay"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/wifi"
	"github.com/charmbracelet/bubbles/spinner"
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
	wifiTable    wifi.TableModel
	timer        timer.Model
	popup        overlay.Model
	notification overlay.Model
	width        int
	height       int
}

func New() Model {
	w := wifi.NewTableModel(30, 20)
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
	var cmd tea.Cmd
	var upd tea.Model
	switch msg := msg.(type) {
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case spinner.TickMsg:
		upd, cmd = m.wifiTable.Update(msg)
		m.wifiTable = upd.(wifi.TableModel)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(msg)
		m.notification = upd.(overlay.Model)
		return m, cmd
	} else if m.popup.IsActive {
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(overlay.Model)
		return m, cmd
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+q", "esc", "ctrl+c":
				return m, tea.Quit
			case "tab":
				if m.state == wifiView {
					m.state = timerView
				} else {
					m.state = wifiView
				}
				return m, cmd
			}
			upd, cmd = m.wifiTable.Update(msg)
			m.wifiTable = upd.(wifi.TableModel)
			return m, cmd
		case PopupContentMsg:
			cmd = m.showPopup(msg)
			return m, cmd
		case NotificationMsg:
			cmd = m.showNotification(string(msg))
			return m, cmd
		}
		upd, cmd = m.wifiTable.Update(msg)
		m.wifiTable = upd.(wifi.TableModel)
		return m, cmd
	}
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

func (m *Model) showPopup(content tea.Model) tea.Cmd {
	m.popup.IsActive = true
	p, cmd := m.popup.Update(overlay.LoadedContentMsg(content))
	m.popup = p.(overlay.Model)
	return cmd
}

func (m *Model) showNotification(text string) tea.Cmd {
	m.notification.IsActive = true
	n, cmd := m.notification.Update(overlay.LoadedContentMsg(label.New(text)))
	m.notification = n.(overlay.Model)
	return cmd
}

type PopupContentMsg tea.Model

func ShowPopup(content tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg(content)
	}
}

type NotificationMsg string

func ShowNotification(notification string) tea.Cmd {
	return func() tea.Msg {
		return NotificationMsg(notification)
	}
}
