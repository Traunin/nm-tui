package ui

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiTableStoredModel struct {
	dataTable table.Model
}

func NewWifiTableStoredTable(width, height int) *WifiTableStoredModel {
	offset := 4
	connectionFlagWidth := 1
	ssidWidth := width - offset - connectionFlagWidth
	cols := []table.Column{
		{Title: "󱘖", Width: connectionFlagWidth},
		{Title: "SSID", Width: ssidWidth},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithWidth(width),
		table.WithHeight(height),
	)
	t.SetStyles(styles.TableStyle)
	m := &WifiTableStoredModel{
		dataTable: t,
	}
	return m
}

func (m WifiTableStoredModel) Init() tea.Cmd {
	return UpdateWifiStoredRows()
}

type storedRowsMsg []table.Row

func (m WifiTableStoredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, UpdateWifiStoredRows()
			// row := m.dataTable.SelectedRow()
			// if row != nil {
			// 	connector := NewWifiConnector(row[1])
			// 	return m, tea.Sequence(SetPopupActivity(true), SetPopupContent(connector))
			// }
			// return m, nil
		}
	case storedRowsMsg:
		m.dataTable.SetRows(msg)
		return m, nil
	}

	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m WifiTableStoredModel) View() string {
	out := m.dataTable.View()

	sb := strings.Builder{}
	sb.WriteString(out)
	sb.WriteString("\n")
	return sb.String()
}

func UpdateWifiStoredRows() tea.Cmd {
	return func() tea.Msg {
		list, err := nmcli.WifiStoredConnections()
		if err != nil {
			logger.Errln(fmt.Errorf("error: %s", err.Error()))
		}
		rows := []table.Row{}
		for _, wifiStored := range list {
			var connectionFlag string
			if wifiStored.Active {
				connectionFlag = ""
			}
			rows = append(rows, table.Row{connectionFlag, wifiStored.Name})
		}
		return storedRowsMsg(rows)
	}
}
