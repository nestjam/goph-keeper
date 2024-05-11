package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nestjam/goph-keeper/internal/tui/auth"
)

var (
	BuildVersion string
	BuildDate    string
)

func main() {
	m := auth.NewLoginModel()
	m.BuildDate = BuildDate
	m.BuildVersion = BuildVersion

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
