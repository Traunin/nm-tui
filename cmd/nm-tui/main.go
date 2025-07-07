package main

import (
	"log"

	"github.com/alphameo/nm-tui/internal/ui/wifi"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(wifi.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
