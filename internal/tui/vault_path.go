package tui

import (
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type VaultPathModel struct {
	input textinput.Model
}

func NewVaultPathModel() VaultPathModel {
	ti := textinput.New()
	ti.Placeholder = "Vault path (e.g. ~/myvault)"
	ti.CharLimit = 256
	ti.Width = 40
	return VaultPathModel{input: ti}
}

func (m *VaultPathModel) Reset() {
	m.input.SetValue("")
	m.input.Focus()
}

func (m VaultPathModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m VaultPathModel) Update(msg tea.Msg, mode VaultMode) (VaultPathModel, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			path := m.input.Value()
			return m, func() tea.Msg {
				return ScreenMsg{
					NextScreen: ScreenPassword,
					Payload: VaultPathResult{
						Mode: mode,
						Path: path,
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

func (m VaultPathModel) View(mode VaultMode) string {
	title := "Open vault – enter path:"
	if mode == VaultModeCreate {
		title = "Create vault – enter path:"
	}
	s := "\n" + title + "\n\n"
	s += m.input.View()
	s += "\n\nEnter to continue, Esc to cancel.\n"
	return s
}
