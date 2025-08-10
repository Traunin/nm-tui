package ui

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiConnectorModel struct {
	SSID     string
	password textinput.Model
	err      error
}

func NewWifiConnector(ssid string) *WifiConnectorModel {
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
	return &WifiConnectorModel{SSID: ssid, password: p}
}

func (m WifiConnectorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WifiConnectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			pw := m.password.Value()
			return m, tea.Sequence(SetPopupActivity(false), WifiConnect(m.SSID, pw))
		case tea.KeyCtrlR:
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
		}
	}

	var cmd tea.Cmd
	m.password, cmd = m.password.Update(msg)
	return m, cmd
}

func (m WifiConnectorModel) View() string {
	inputField := styles.BorderStyle.Render(m.password.View())
	return fmt.Sprintf("SSID: %s\n%v", m.SSID, inputField)
}
