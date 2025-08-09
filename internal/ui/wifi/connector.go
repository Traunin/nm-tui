// Package wifi contains a bunch of windows for wifi control
package wifi

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/overlay"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectorModel struct {
	SSID     string
	password textinput.Model
	err      error
}

type errMsg error

func NewConnector(ssid string) *ConnectorModel {
	p := textinput.New()
	p.Focus()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "Password"
	pw, err := nmcli.WifiGetPassword(&ssid)
	if err == nil {
		p.SetValue(pw)
	}
	return &ConnectorModel{SSID: ssid, password: p}
}

func (m ConnectorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConnectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			go func() {
				pw := m.password.Value()
				nmcli.WifiConnect(&m.SSID, &pw)
			}()
			return m, overlay.Close()
		case tea.KeyCtrlR:
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.password, cmd = m.password.Update(msg)
	return m, cmd
}

func (m ConnectorModel) View() string {
	pwInput := styles.BorderStyle.Render(m.password.View())
	return fmt.Sprintf("SSID: %s\n%v", m.SSID, pwInput)
}
