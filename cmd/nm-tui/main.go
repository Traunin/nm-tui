package main

import (
	"log"

	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.Init("./log")
	logger.InfoLog.Println("The program is running")
	defer logger.InfoLog.Println("Program is closed")
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
