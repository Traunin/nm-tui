package connections

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiStoredInfoModel struct {
	SSID string
}

func NewStoredInfoModel(ssid string) *WifiStoredInfoModel {
	return &WifiStoredInfoModel{SSID: ssid}
}

func (m WifiStoredInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WifiStoredInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m, nil
		}
	}
	return m, nil
}

func (m WifiStoredInfoModel) View() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n", m.SSID)
	return sb.String()
}
