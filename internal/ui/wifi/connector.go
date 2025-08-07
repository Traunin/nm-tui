package wifi

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/logger"
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

		// case tea.KeyCtrlQ, tea.KeyEsc:
		// 	return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.password, cmd = m.password.Update(msg)
	logger.InfoLog.Println("l")
	return m, cmd
}

func (m ConnectorModel) View() string {
	return fmt.Sprintf("SSID: %s\n%s", m.SSID, styles.BorderStyle.Render(m.password.View()))
}
