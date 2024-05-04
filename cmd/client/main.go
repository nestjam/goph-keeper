package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nestjam/goph-keeper/internal/tui/auth"
)

func main() {
	m := auth.NewLoginModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
