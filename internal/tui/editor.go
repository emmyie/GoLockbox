package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorModel struct {
	textarea textarea.Model
	index    int
	filename string
}

func NewEditorModel() EditorModel {
	ta := textarea.New()
	ta.SetWidth(80)
	ta.SetHeight(20)
	ta.ShowLineNumbers = true

	return EditorModel{
		textarea: ta,
		index:    -1,
	}
}

func (m *EditorModel) SetFile(index int, name string, data []byte) {
	m.index = index
	m.filename = name
	m.textarea.SetValue(string(data))
	m.textarea.Focus()
}

func (m EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+s":
			idx := m.index
			data := []byte(m.textarea.Value())
			name := m.filename

			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload: SaveFileMsg{
						Index: idx,
						Data:  data,
						Name:  name,
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

func (m EditorModel) View() string {
	s := "\nEdit file – Ctrl+S to save, Esc to cancel\n\n"
	s += fmt.Sprintf("Filename: %s\n\n", m.filename)
	s += m.textarea.View()
	s += "\n"
	return s
}
