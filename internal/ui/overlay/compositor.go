package overlay

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func Compose(fg, bg string, x, y int) string {
	fgW, fgH := lipgloss.Size(fg)
	bgW, bgH := lipgloss.Size(bg)
	fgXmax := x + fgW
	fgYmax := y + fgH

	if (fgW >= bgW && fgH >= bgH) || x >= bgW || y >= bgH || fgXmax < 0 || fgYmax < 0 {
		return bg
	}

	fgLines := lines(fg)
	bgLines := lines(bg)

	var sb strings.Builder

	var fgInd int
	if y < 0 {
		fgInd -= y
	}

	for bgY, bgLine := range bgLines {
		if bgY > 0 {
			sb.WriteByte('\n')
		}
		if bgY < y || bgY >= fgYmax {
			sb.WriteString(bgLine)
			continue
		}
		if x > 0 {
			left := ansi.Truncate(bgLine, x, "")
			sb.WriteString(left)
		}
		fgLine := fgLines[fgInd]
		fgInd++

		if x < 0 {
			fgLine = ansi.TruncateLeft(fgLine, -x, "")
		}

		if fgXmax <= bgW {
			sb.WriteString(fgLine)
		} else {
			sb.WriteString(ansi.Truncate(fgLine, fgW-fgXmax+bgW, ""))
			continue
		}
		right := ansi.TruncateLeft(bgLine, fgXmax, "")
		sb.WriteString(right)
	}
	return sb.String()
}

func lines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}
