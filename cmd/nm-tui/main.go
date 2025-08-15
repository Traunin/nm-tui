package main

import (
	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.Init("./log")
	logger.Informln("The program is running")
	defer logger.Informln("Program is closed")
	logger.Level = logger.ErrorsLvl
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Errln(err)
	}
}
