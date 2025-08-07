// Package session contains Model, which represents main window of TUI
package session

import (
	"fmt"
	"time"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/overlay"
	"github.com/alphameo/nm-tui/internal/ui/styles"
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
	state        sessionState
	wifi         wifi.Model
	timer        timer.Model
	popup        overlay.Model
	notification overlay.Model
	width        int
	height       int
}

func New() Model {
	w := wifi.New(30, 20)
	t := timer.New(time.Hour)
	p := overlay.New(NewTextModel())
	p.XAnchor = overlay.Center
	p.YAnchor = overlay.Center
	n := overlay.New(NewTextModel())
	n.XAnchor = overlay.Center
	n.YAnchor = overlay.Center
	m := Model{
		wifi:         *w,
		timer:        t,
		popup:        *p,
		notification: *n,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		m.wifi.Init(),
		m.popup.Init(),
		m.notification.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	var upd tea.Model
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(msg)
		m.notification = upd.(overlay.Model)
		cmds = append(cmds, cmd)
	} else if m.popup.IsActive {
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(overlay.Model)
		cmds = append(cmds, cmd)
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
			case "o":
				cmd = m.showPopup(nil)
				cmds = append(cmds, cmd)
			case "n":
				m.notify("xddddddd\nddddd")
			}
		}
	}
	size, ok := msg.(tea.WindowSizeMsg)
	if ok {
		m.width = size.Width
		m.height = size.Height
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
	if content != nil {
		m.popup.Content = content
	}
	return m.popup.Content.Init()
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
