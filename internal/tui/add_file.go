package tui

import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AddFileModel struct {
	input textinput.Model
}

func NewAddFileModel() AddFileModel {
	ti := textinput.New()
	ti.Placeholder = "Path to source file"
	ti.CharLimit = 512
	ti.Width = 60
	return AddFileModel{input: ti}
}

func (m *AddFileModel) Init() tea.Cmd {
	m.input.SetValue("")
	m.input.Focus()
	return textinput.Blink
}

func (m AddFileModel) Update(msg tea.Msg) (AddFileModel, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			src := m.input.Value()
			if src == "" {
				return m, nil
			}
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload: AddFileMsg{
						SrcPath: src,
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

func (m AddFileModel) View() string {
	s := "\nAdd existing file – enter source path\n\n"
	s += m.input.View()
	s += "\n\nEnter to add, Esc to cancel.\n"
	return s
}
