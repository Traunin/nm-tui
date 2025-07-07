package ui

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type updatedRowsMsg []table.Row

type WifiModel struct {
	wifiTable       table.Model
	updatingSpinner spinner.Model
	updating        bool
}

func NewModel() WifiModel {
	columns := []table.Column{
		{Title: "SSID", Width: 16},
		{Title: "Signal", Width: 8},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := spinner.New()
	m := WifiModel{wifiTable: t, updatingSpinner: s, updating: true}
	return m
}

func (m WifiModel) Init() tea.Cmd {
	return tea.Batch(m.updatingSpinner.Tick, m.updateWifiList())
}

func (m WifiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "r":
			if m.updating {
				return m, nil
			}
			m.updating = true
			cmds := []tea.Cmd{
				m.updateWifiList(),     // start loading
				m.updatingSpinner.Tick, // restart spinner animation
			}
			return m, tea.Batch(cmds...)
		}
	case updatedRowsMsg:
		m.updating = false
		m.wifiTable.SetRows(msg)
		return m, nil
	}
	if m.updating {
		m.updatingSpinner, cmd = m.updatingSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.wifiTable, cmd = m.wifiTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m WifiModel) View() string {
	out := m.wifiTable.View()
	if m.updating {
		out += "\n" + m.updatingSpinner.View()
	}
	return out
}

func (m *WifiModel) updateWifiList() tea.Cmd {
	return func() tea.Msg {
		rows := getWifiRows()
		return updatedRowsMsg(rows)
	}
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
