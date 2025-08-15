package ui

import (
	"strings"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WifiTableModel struct {
	tables    []tea.Model
	tabTitles []string
	activeTab int
}

func NewWifiTableModel(width, height int) *WifiTableModel {
	current := NewWifiTableCurrentTable(width, height)
	stored := NewWifiTableStoredTable(width, height)
	ts := []tea.Model{current, stored}
	tabTitles := &[]string{"Current", "Stored"}
	m := &WifiTableModel{
		tables:    ts,
		tabTitles: *tabTitles,
		activeTab: 0,
	}
	return m
}

func (m WifiTableModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tables {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func (m WifiTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "]", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabTitles)-1)
			return m, m.tables[m.activeTab].Init()
		case "[", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, m.tables[m.activeTab].Init()
		}
	}

	var cmd tea.Cmd
	m.tables[m.activeTab], cmd = m.tables[m.activeTab].Update(msg)
	return m, cmd
}

func (m WifiTableModel) View() string {
	out := m.tables[m.activeTab].View()

	fullWidth := lipgloss.Width(out) + 2
	tabCount := len(m.tabTitles)
	tabWidth := fullWidth/tabCount - 2
	tail := fullWidth % tabCount
	logger.Debugln(fullWidth, tabCount, tabWidth, tail)
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
		logger.Debugln(tabWidth, lipgloss.Width(tabView))
		renderedTabs = append(renderedTabs, tabView)
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

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
	return sb.String()
}

type updatedRowsMsg []table.Row
