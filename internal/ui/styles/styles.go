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

func ConstructTabBar(
	titles []string,
	activeStyle,
	inactiveStyle lipgloss.Style,
	fullWidth int,
	active int,
) string {
	tabCount := len(titles)
	tabWidth := fullWidth/tabCount - 2
	tail := fullWidth % tabCount
	var renderedTabs []string
	for i, t := range titles {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(titles)-1, i == active
		if isActive {
			style = activeStyle
		} else {
			style = inactiveStyle
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

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
