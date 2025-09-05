package connections

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/controls"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WifiConnectorModel struct {
	ssid     string
	password textinput.Model
}

func NewWifiConnector() *WifiConnectorModel {
	p := textinput.New()
	p.Focus()
	p.Width = 20
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.Placeholder = "Password"
	return &WifiConnectorModel{password: p}
}

func (m *WifiConnectorModel) setNew(ssid string) {
	m.ssid = ssid
	pw, err := nmcli.WifiGetPassword(ssid)
	if err == nil {
		m.password.SetValue(pw)
	}
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
			return m, tea.Sequence(
				controls.SetPopupActivity(false),
				WifiConnect(m.ssid, pw),
			)
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
	inputField := lipgloss.
		NewStyle().
		BorderStyle(styles.BorderStyle).
		Render(m.password.View())
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "SSID: %s\n%v", m.ssid, inputField)
	return sb.String()
}
