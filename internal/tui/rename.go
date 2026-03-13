package tui

import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type RenameModel struct {
	input textinput.Model
	index int
}

func NewRenameModel() RenameModel {
	ti := textinput.New()
	ti.Placeholder = "New filename"
	ti.CharLimit = 256
	ti.Width = 40
	return RenameModel{
		input: ti,
		index: -1,
	}
}

func (m *RenameModel) SetTarget(index int, currentName string) {
	m.index = index
	m.input.SetValue(currentName)
	m.input.Focus()
}

func (m RenameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RenameModel) Update(msg tea.Msg) (RenameModel, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			newName := m.input.Value()
			idx := m.index
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload: RenameResultMsg{
						Index:   idx,
						NewName: newName,
					},
				}
			}
		case "esc":
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload:    BackToBrowserMsg{},
				}
			}
		}
	}

	return m, cmd
}

func (m RenameModel) View() string {
	s := "\nRename file – enter new name\n\n"
	s += m.input.View()
	s += "\n\nEnter to rename, Esc to cancel.\n"
	return s
}
