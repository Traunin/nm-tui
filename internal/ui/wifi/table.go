package wifi

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/overlay"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type updatedRowsMsg []table.Row

type TableModel struct {
	wifiTable       table.Model
	updatingSpinner spinner.Model
	updating        bool
}

func NewTableModel(width int, height int) *TableModel {
	offset := 4
	signalWidth := 3
	ssidWidth := width - signalWidth - offset
	cols := []table.Column{
		{Title: "SSID", Width: ssidWidth},
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
	m := &TableModel{wifiTable: t, updatingSpinner: s, updating: true}
	return m
}

func (m TableModel) Init() tea.Cmd {
	return tea.Batch(m.updatingSpinner.Tick, m.updateWifiList())
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			cmds := []tea.Cmd{
				m.updateWifiList(),
				m.updatingSpinner.Tick,
			}
			return m, tea.Batch(cmds...)
		case "enter":
			row := m.wifiTable.SelectedRow()
			if row != nil {
				return m, overlay.LoadContent(NewConnector(row[0]))
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

func (m TableModel) View() string {
	out := m.wifiTable.View()
	var symbol string
	if m.updating {
		symbol = m.updatingSpinner.View()
	} else {
		symbol = "󰄬 "
	}
	out += "\n" + lipgloss.Place(m.wifiTable.Width(), 1, lipgloss.Center, lipgloss.Center, symbol)
	return styles.BorderStyle.Render(out)
}

func (m *TableModel) updateWifiList() tea.Cmd {
	return func() tea.Msg {
		rows := getWifiRows()
		return updatedRowsMsg(rows)
	}
}

func getWifiRows() []table.Row {
	list, err := nmcli.WifiScan()
	if err != nil {
		fmt.Println(fmt.Errorf("error: %s", err.Error()))
	}
	rows := []table.Row{}
	for _, wifiNet := range list {
		rows = append(rows, table.Row{wifiNet.SSID, fmt.Sprint(wifiNet.Signal)})
	}
	return rows
}
