package main

import (
	"log"

	"github.com/alphameo/nm-tui/internal/ui/session"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(session.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
