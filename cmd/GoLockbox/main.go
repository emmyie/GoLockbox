package main

import (
	"fmt"
	"os"

	"GoLockbox/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(tui.NewRootModel()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
