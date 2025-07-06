package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/alphameo/nm-tui/internal/nmcli"
	tea "github.com/charmbracelet/bubbletea"
)

type Model int

type tickMsg time.Time

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick()
	}
	return m, nil
}

func (m Model) View() string {
	WifiList, err := nmcli.ScanWifi()
	if err != nil {
		return err.Error()
	}
	sb := strings.Builder{}
	for i, wifiNet := range WifiList {
		line := fmt.Sprintf("%v: %s\t%v\n", i+1, wifiNet.SSID, wifiNet.Signal)
		sb.WriteString(line)
	}
	return sb.String()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
