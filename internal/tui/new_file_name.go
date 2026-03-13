package tui

import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type NewFileNameModel struct {
	input textinput.Model
}

func NewNewFileNameModel() NewFileNameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter new filename"
	ti.CharLimit = 256
	ti.Width = 40
	return NewFileNameModel{input: ti}
}

func (m *NewFileNameModel) Init() tea.Cmd {
	m.input.SetValue("")
	m.input.Focus()
	return textinput.Blink
}

func (m NewFileNameModel) Update(msg tea.Msg) (NewFileNameModel, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			name := m.input.Value()
			if name == "" {
				return m, nil
			}
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenFileEditor,
					Payload: NewFileNameResult{
						Name: name,
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

func (m NewFileNameModel) View() string {
	s := "\nCreate new file – enter filename\n\n"
	s += m.input.View()
	s += "\n\nEnter to continue, Esc to cancel.\n"
	return s
}
