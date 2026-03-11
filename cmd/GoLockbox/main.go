package main

import (
	"log"

	"GoLockbox/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(tui.InitialModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
