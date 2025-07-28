// Package wifi provides interaction with wifi networks from nmcli
package wifi

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type updatedRowsMsg []table.Row

type Model struct {
	wifiTable       table.Model
	updatingSpinner spinner.Model
	updating        bool
}

func New(width int, height int) *Model {
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
	m := &Model{wifiTable: t, updatingSpinner: s, updating: true}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.updatingSpinner.Tick, m.updateWifiList())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
	out := m.wifiTable.View()
	if m.updating {
		out += "\n" + m.updatingSpinner.View()
	} else {
		out += "\n󰄬 "
	}
	return styles.BorderStyle.Render(out)
}

func (m *Model) updateWifiList() tea.Cmd {
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
