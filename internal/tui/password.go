package tui

import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PasswordModel struct {
	input textinput.Model
}

func NewPasswordModel() PasswordModel {
	ti := textinput.New()
	ti.Placeholder = "Password"
	ti.CharLimit = 128
	ti.Width = 40
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'
	return PasswordModel{input: ti}
}

func (m *PasswordModel) Reset() {
	m.input.SetValue("")
	m.input.Focus()
}

func (m PasswordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PasswordModel) Update(msg tea.Msg, mode VaultMode, path string) (PasswordModel, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			pass := m.input.Value()
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenBrowser,
					Payload: PasswordResult{
						Mode:     mode,
						Path:     path,
						Password: pass,
					},
				}
			}
		case "esc":
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenMenu,
					Payload:    nil,
				}
			}
		}
	}

	return m, cmd
}

func (m PasswordModel) View(mode VaultMode, path string) string {
	title := "Open vault – enter password:"
	if mode == VaultModeCreate {
		title = "Create vault – enter password:"
	}
	s := "\n" + title + "\n\n"
	s += m.input.View()
	s += "\n\nEnter to continue, Esc to cancel.\n"
	return s
}
