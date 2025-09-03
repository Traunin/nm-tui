package connections

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/controls"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type WifiStoredModel struct {
	dataTable table.Model
}

func NewWifiStored(width, height int) *WifiStoredModel {
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
	m := &WifiStoredModel{
		dataTable: t,
	}
	return m
}

func (m WifiStoredModel) Init() tea.Cmd {
	return UpdateWifiStoredRows()
}

type storedRowsMsg []table.Row

func (m WifiStoredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				connector := NewStoredInfoModel(row[1])
				return m, tea.Sequence(controls.SetPopupActivity(true), controls.SetPopupContent(connector))
			}
			return m, nil
		case "r":
			return m, UpdateWifiStoredRows()
		case "d":
			row := m.dataTable.SelectedRow()
			cursor := m.dataTable.Cursor()
			if cursor != 0 {
				m.dataTable.SetCursor(cursor - 1)
			}
			return m, tea.Sequence(DeleteConnection(row[1]), UpdateWifiStoredRows())
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

func (m WifiStoredModel) View() string {
	view := m.dataTable.View()

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n", view)
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
