package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	choices []string
	cursor  int
}

func NewMenuModel() MenuModel {
	return MenuModel{
		choices: []string{"New Vault", "Open Vault"},
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Reset() {
	m.cursor = 0
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			choice := MenuChoice(m.cursor)
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenVaultPath,
					Payload:    choice,
				}
			}
		}
	}
	return m, nil
}

func (m MenuModel) View(errMsg string) string {
	s := "\nMain menu:\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	if errMsg != "" {
		s += "\nError: " + errMsg + "\n"
	}
	s += "\nUse ↑/↓ and Enter. Ctrl+C to quit.\n"
	return s
}
