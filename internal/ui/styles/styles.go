// Package styles contains common styles for whole app
package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	BorderStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	TableStyle        = makeTableStyle()
	InactiveTabBorder = makeTabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = makeTabBorderWithBottom("┘", " ", "└")
	InactiveTabStyle  = lipgloss.NewStyle().Border(InactiveTabBorder, true).Padding(0, 1)
	ActiveTabStyle    = InactiveTabStyle.Border(ActiveTabBorder, true)
)

func makeTableStyle() table.Styles {
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return style
}

func makeTabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
