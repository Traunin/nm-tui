package ui

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type updatedRowsMsg []table.Row

type WifiTableModel struct {
	wifiTable       table.Model
	updatingSpinner spinner.Model
	updating        bool
}

func NewWifiTableModel(width int, height int) *WifiTableModel {
	offset := 8
	signalWidth := 3
	connectionFlagWidth := 1
	securityWidth := 10
	ssidWidth := width - signalWidth - offset - connectionFlagWidth - securityWidth
	cols := []table.Column{
		{Title: "󱘖", Width: connectionFlagWidth},
		{Title: "SSID", Width: ssidWidth},
		{Title: "Security", Width: securityWidth},
		{Title: "", Width: signalWidth},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithWidth(width),
		table.WithHeight(height),
	)
	t.SetStyles(styles.TableStyle)
	s := spinner.New()
	m := &WifiTableModel{wifiTable: t, updatingSpinner: s, updating: true}
	return m
}

func (m WifiTableModel) Init() tea.Cmd {
	return tea.Batch(m.updatingSpinner.Tick, UpdateWifiList)
}

func (m WifiTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if m.updating {
				return m, nil
			}
			m.updating = true
			cmds = []tea.Cmd{
				UpdateWifiList,
				m.updatingSpinner.Tick,
			}
			return m, tea.Batch(cmds...)
		case "enter":
			row := m.wifiTable.SelectedRow()
			if row != nil {
				return m, ShowPopup(NewWifiConnector(row[1]))
			}
			return m, nil
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

func (m WifiTableModel) View() string {
	out := m.wifiTable.View()
	var symbol string
	if m.updating {
		symbol = m.updatingSpinner.View()
	} else {
		symbol = "󰄬"
	}
	out += "\n" + lipgloss.Place(m.wifiTable.Width(), 1, lipgloss.Center, lipgloss.Center, symbol)
	return styles.BorderStyle.Render(out)
}

func UpdateWifiList() tea.Msg {
	rows := getWifiRows()
	return updatedRowsMsg(rows)
}

func getWifiRows() []table.Row {
	list, err := nmcli.WifiScan()
	if err != nil {
		fmt.Println(fmt.Errorf("error: %s", err.Error()))
	}
	rows := []table.Row{}
	for _, wifiNet := range list {
		var connectionFlag string
		if wifiNet.Active {
			connectionFlag = ""
		}
		rows = append(rows, table.Row{connectionFlag, wifiNet.SSID, wifiNet.Security, fmt.Sprint(wifiNet.Signal)})
	}
	return rows
}
