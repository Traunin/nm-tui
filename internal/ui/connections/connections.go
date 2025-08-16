// Package connections provides tabbed tables with main iformation about variable connections
package connections

import (
	"strings"

	"github.com/alphameo/nm-tui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabTables []tea.Model
	tabTitles []string
	activeTab int
}

func New(width, height int) *Model {
	wifiAvailable := NewWifiAvailable(width, height)
	wifiStored := NewWifiStored(width, height)
	tabTables := []tea.Model{wifiAvailable, wifiStored}
	tabTitles := &[]string{"Current", "Stored"}
	m := &Model{
		tabTables: tabTables,
		tabTitles: *tabTitles,
		activeTab: 0,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tabTables {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "]", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabTitles)-1)
			return m, m.tabTables[m.activeTab].Init()
		case "[", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, m.tabTables[m.activeTab].Init()
		}
	}

	var cmd tea.Cmd
	m.tabTables[m.activeTab], cmd = m.tabTables[m.activeTab].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	view := m.tabTables[m.activeTab].View()

	fullWidth := lipgloss.Width(view) + 2
	tabCount := len(m.tabTitles)
	tabWidth := fullWidth/tabCount - 2
	tail := fullWidth % tabCount
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
	sb.WriteString(style.Render(view))
	return sb.String()
}
