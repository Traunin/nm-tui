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
	dataTable  table.Model
	storedInfo WifiStoredInfoModel
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
	s := NewStoredInfoModel()
	m := &WifiStoredModel{
		dataTable:  t,
		storedInfo: *s,
	}
	return m
}

func (m WifiStoredModel) Init() tea.Cmd {
	return m.UpdateRows()
}

type storedRowsMsg []table.Row

func (m WifiStoredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				m.storedInfo.setNew(row[1])
				return m, tea.Sequence(controls.SetPopupActivity(true), controls.SetPopupContent(m.storedInfo))
			}
			return m, nil
		case "r":
			return m, m.UpdateRows()
		case "d":
			row := m.dataTable.SelectedRow()
			cursor := m.dataTable.Cursor()
			if cursor == len(m.dataTable.Rows())-1 {
				m.dataTable.SetCursor(cursor - 1)
			}
			return m, tea.Sequence(controls.DeleteConnection(row[1]), m.UpdateRows())
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

func (m WifiStoredModel) UpdateRows() tea.Cmd {
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
