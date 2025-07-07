package ui

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	wifiTable table.Model
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "r":
			rows := getWifiRows()
			m.wifiTable.SetRows(rows)
			m.wifiTable, cmd = m.wifiTable.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func getWifiRows() []table.Row {
	list, err := nmcli.ScanWifi()
	if err != nil {
		fmt.Println(fmt.Errorf("error: %s", err.Error()))
	}
	rows := []table.Row{}
	for _, wifiNet := range list {
		rows = append(rows, table.Row{wifiNet.SSID, fmt.Sprint(wifiNet.Signal)})
	}
	return rows
}

func (m Model) View() string {
	return m.wifiTable.View()
}

func InitialModel() Model {
	columns := []table.Column{
		{Title: "SSID", Width: 16},
		{Title: "Signal", Width: 8},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(getWifiRows()),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	m := Model{wifiTable: t}
	return m
}
