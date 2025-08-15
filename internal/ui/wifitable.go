package ui

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wifiState int

const (
	Scanning wifiState = iota
	Connecting
	None
)

func (s *wifiState) String() string {
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
	indicatorState   wifiState
	tabTitles        []string
	activeTab        int
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
	tabTitles := []string{"Current", "Stored"}

	m := &WifiTableModel{
		dataTable:        t,
		indicatorSpinner: s,
		indicatorState:   Scanning,
		tabTitles:        tabTitles,
	}
	return m
}

func (m WifiTableModel) Init() tea.Cmd {
	return tea.Batch(UpdateWifiRows())
}

func (m WifiTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if m.indicatorState != None {
				return m, nil
			}
			return m, UpdateWifiRows()
		case "enter":
			row := m.dataTable.SelectedRow()
			if row != nil {
				connector := NewWifiConnector(row[1])
				return m, tea.Sequence(SetPopupActivity(true), SetPopupContent(connector))
			}
			return m, nil
		case "]", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabTitles)-1)
			return m, nil
		case "[", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	case updatedRowsMsg:
		m.dataTable.SetRows(msg)
		return m, nil
	case WifiIndicatorStateMsg:
		m.indicatorState = wifiState(msg)
		if m.indicatorState == None {
			return m, nil
		}
		return m, m.indicatorSpinner.Tick
	case AfterWifiConnectionMsg:
		return m, tea.Cmd(msg)
	}

	var cmd tea.Cmd
	if m.indicatorState != None {
		m.indicatorSpinner, cmd = m.indicatorSpinner.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}
	m.dataTable, cmd = m.dataTable.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m WifiTableModel) View() string {
	out := m.dataTable.View()

	fullWidth := lipgloss.Width(out) + 2
	tabCount := len(m.tabTitles)
	tabWidth := fullWidth/tabCount - 2
	tail := fullWidth % tabCount
	logger.InfoLog.Println(fullWidth, tabCount, tabWidth, tail)
	var renderedTabs []string
	for i, t := range m.tabTitles {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabTitles)-1, i == m.activeTab
		if isActive {
			style = styles.ActiveTabStyle
		} else {
			style = styles.InactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		if tail > 0 {
			style = style.Width(tabWidth + 1)
			tail--
		} else {
			style = style.Width(tabWidth)
		}
		tabView := style.Render(t)
		logger.InfoLog.Println(tabWidth, lipgloss.Width(tabView))
		renderedTabs = append(renderedTabs, tabView)
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	var symbol string
	if m.indicatorState != None {
		symbol = fmt.Sprintf("%s %s", m.indicatorState.String(), m.indicatorSpinner.View())
	} else {
		symbol = "󰄬"
	}

	statusline := lipgloss.Place(m.dataTable.Width(), 1, lipgloss.Center, lipgloss.Center, symbol)

	sb := strings.Builder{}
	sb.WriteString(tabRow)
	sb.WriteString("\n")
	borderStyle := styles.BorderStyle.GetBorderStyle()
	borderStyle.Top = ""
	borderStyle.TopLeft = "│"
	borderStyle.TopRight = "│"
	var style lipgloss.Style
	style = style.Border(borderStyle)
	sb.WriteString(style.Render(out))
	sb.WriteString("\n")
	sb.WriteString(statusline)
	return sb.String()
}

type updatedRowsMsg []table.Row

func UpdateWifiRows() tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorState(Scanning),
		func() tea.Msg {
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
			return updatedRowsMsg(rows)
		},
		SetWifiIndicatorState(None))
}

type WifiIndicatorStateMsg wifiState

func SetWifiIndicatorState(state wifiState) tea.Cmd {
	return func() tea.Msg {
		return WifiIndicatorStateMsg(state)
	}
}

type AfterWifiConnectionMsg tea.Cmd

func WifiConnect(ssid, password string) tea.Cmd {
	return tea.Sequence(
		SetWifiIndicatorState(Connecting),
		func() tea.Msg {
			err := nmcli.WifiConnect(ssid, password)
			if err == nil {
				return AfterWifiConnectionMsg(UpdateWifiRows())
			} else {
				error := fmt.Sprintf("Connection interrupted: %s", err.Error())
				return AfterWifiConnectionMsg(tea.Sequence(SetNotificationActivity(true), SetNotificationText(error)))
			}
		},
		SetWifiIndicatorState(None))
}
