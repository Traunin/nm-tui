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

type tableSpinnerState int

const (
	Scanning tableSpinnerState = iota
	Connecting
	None
)

func (s *tableSpinnerState) String() string {
	switch *s {
	case Scanning:
		return "Scanning"
	case Connecting:
		return "Connecting"
	case None:
		return ""
	default:
		return "Undefined!!!"
	}
}

type WifiTableModel struct {
	dataTable        table.Model
	indicatorSpinner spinner.Model
	indicatorState   tableSpinnerState
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
	m := &WifiTableModel{dataTable: t, indicatorSpinner: s, indicatorState: Scanning}
	return m
}

func (m WifiTableModel) Init() tea.Cmd {
	return tea.Batch(m.indicatorSpinner.Tick, UpdateWifiRows)
}

func (m WifiTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if m.indicatorState != None {
				return m, nil
			}
			m.indicatorState = Scanning
			return m, tea.Batch(UpdateWifiRows, m.indicatorSpinner.Tick)
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				return m, ShowPopup(NewWifiConnector(row[1]))
			}
			return m, nil
		}
	case updatedRowsMsg:
		m.indicatorState = None
		m.dataTable.SetRows(msg)
		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.indicatorState != None {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.dataTable, cmd = m.dataTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m WifiTableModel) View() string {
	out := m.dataTable.View()
	var symbol string
	if m.indicatorState != None {
		symbol = fmt.Sprintf("%s %s", m.indicatorState.String(), m.indicatorSpinner.View())
	} else {
		symbol = "󰄬"
	}
	out += "\n" + lipgloss.Place(m.dataTable.Width(), 1, lipgloss.Center, lipgloss.Center, symbol)
	return styles.BorderStyle.Render(out)
}

type updatedRowsMsg []table.Row

func UpdateWifiRows() tea.Msg {
	rows := getWifiRows()
	return updatedRowsMsg(rows)
}

func getWifiRows() []table.Row {
	list, err := nmcli.WifiScan()
	if err != nil {
		logger.ErrorLog.Println(fmt.Errorf("error: %s", err.Error()))
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

type tableSpinnerStateMsg tableSpinnerState

func SetTableSpinnerState(state tableSpinnerState) tea.Cmd {
	return func() tea.Msg {
		return tableSpinnerState(state)
	}
}

type wifiConnectionMsg struct {
	SSID     string
	password string
}

func WifiConnect(ssid, password string) tea.Cmd {
	return func() tea.Msg {
		return wifiConnectionMsg{SSID: ssid, password: password}
	}
}
